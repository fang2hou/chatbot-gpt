package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"

	"chatbot-gpt/internal/locale"
)

// predictTokens predicts the number of tokens usage for the given message.
func predictTokens(messages []openai.ChatCompletionMessage, includeAssistantSignal bool) int {
	numTokens := 0

	if includeAssistantSignal {
		numTokens += 3
	}

	tokensPerMessage := 0
	tokensPerName := 0

	switch Model.ID {
	case "gpt-3.5-turbo":
	case "gpt-3.5-turbo-0301":
		tokensPerMessage = 4
		tokensPerName = -1
	case "gpt-4":
	case "gpt-4-0314":
		tokensPerMessage = 3
		tokensPerName = 1
	case "gpt-4-32k":
	case "gpt-4-32k-0314":
		tokensPerMessage = 3
		tokensPerName = 1
	}

	for _, message := range messages {
		numTokens += tokensPerMessage
		if message.Name != "" {
			numTokens += tokensPerName
		}

		numTokens += len(TokenPredictionModel.Encode(message.Role, nil, nil))
		numTokens += len(TokenPredictionModel.Encode(message.Content, nil, nil))
	}

	return numTokens
}

// getTokenCostPriceString returns the cost price of the given number of tokens.
func getTokenCostPriceString(numTokens int) string {
	numDollars := float64(numTokens) * 0.002 / 1000
	numYen := numDollars * 132.45
	numYuan := numDollars * 6.88

	return fmt.Sprintf(
		"💠 %d  →  🇺🇸 $%.3f / 🇯🇵 ￥%.3f / 🇨🇳 ￥%.3f",
		numTokens, numDollars, numYen, numYuan,
	)
}

// storeInteraction stores the interaction between the user and the assistant.
func storeInteraction(
	userID string, userMessage *openai.ChatCompletionMessage, numUserMessageToken int,
	assistantMessage *openai.ChatCompletionMessage, numAssistantMessageToken int,
) error {
	if err := MessageDatabase.Store(userID, userMessage, numUserMessageToken); err != nil {
		Logger.Debug("failed to store response message", zap.Error(err))
		return err
	}

	if err := MessageDatabase.Store(userID, assistantMessage, numAssistantMessageToken); err != nil {
		Logger.Debug("failed to store response message", zap.Error(err))
		return err
	}

	return nil
}

func sendDiscordResponseWithStream(
	stream *openai.ChatCompletionStream, interval time.Duration,
	s *discordgo.Session, guildID, channelID, messageID string,
	lang locale.Language, numPromptTokens int,
) (*openai.ChatCompletionMessage, int, error) {
	var currentResponse *discordgo.Message
	var currentResponseString string
	var allResponseString string

	if resp, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: Localizer.Fetch("wait_for_response", lang),
		Reference: &discordgo.MessageReference{
			MessageID: messageID,
			GuildID:   guildID,
		},
	}); err != nil {
		return nil, 0, err
	} else {
		currentResponse = resp
	}

	lastSentTime := time.Now()

	tryUpdateResponse := func() error {
		if len(currentResponseString) > 2000 {
			firstPart := ""

			for i := 1950; i > 0; i-- {
				if currentResponseString[i] == '\n' {
					firstPart = currentResponseString[:i]
					currentResponseString = currentResponseString[i+1:]
					break
				}
			}

			if _, err := s.ChannelMessageEdit(currentResponse.ChannelID, currentResponse.ID, firstPart); err != nil {
				return err
			}

			if newMessage, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
				Content: currentResponseString,
				Reference: &discordgo.MessageReference{
					MessageID: messageID,
					GuildID:   guildID,
				},
			}); err != nil {
				return err
			} else {
				currentResponse = newMessage
			}
		} else {
			if _, err := s.ChannelMessageEdit(channelID, currentResponse.ID, currentResponseString); err != nil {
				return err
			}
		}

		lastSentTime = time.Now()

		return nil
	}

	for {
		resp, streamErr := stream.Recv()
		if errors.Is(streamErr, io.EOF) || resp.Choices[0].FinishReason == "stop" {
			break
		}

		if len(resp.Choices) > 0 {
			content := resp.Choices[0].Delta.Content
			if len(content) > 0 {
				currentResponseString += content
				allResponseString += content
			}
		}

		if time.Since(lastSentTime) > interval {
			if len(currentResponseString) > 0 {
				if err := tryUpdateResponse(); err != nil {
					return nil, 0, err
				}
			}
		}
	}

	message := &openai.ChatCompletionMessage{
		Content: allResponseString,
		Role:    openai.ChatMessageRoleAssistant,
	}

	numToken := predictTokens([]openai.ChatCompletionMessage{*message}, false)
	currentResponseString += "\n\n" + getTokenCostPriceString(numToken+numPromptTokens)

	if err := tryUpdateResponse(); err != nil {
		return nil, 0, err
	}

	return message, numToken, nil
}

// sendErrorMessage sends an error message.
func sendErrorMessage(s *discordgo.Session, data *discordgo.MessageCreate, lang locale.Language, errorMessage string) {
	if _, err := s.ChannelMessageSendComplex(data.ChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       Localizer.Fetch("error", lang),
				Description: Localizer.Fetch(errorMessage, lang),
				Timestamp:   time.Now().Format(time.RFC3339),
				Color:       0xCC0000,
			},
		},
	}); err != nil {
		Logger.Debug("failed to send message", zap.Error(err))
	}
}

// chatChanel handles the chat channel.
func chatChanel(s *discordgo.Session, data *discordgo.MessageCreate) bool {
	// Only respond to messages that start with the prefix
	if data.Author.ID == s.State.User.ID {
		return false
	}

	serverConfig, sConfigOk := ServerConfigMap[data.GuildID]
	if !sConfigOk {
		return false
	}

	channelConfig, cConfigOk := serverConfig.ChatChannels[data.ChannelID]
	if !cConfigOk {
		return false
	}

	if err := s.ChannelTyping(data.ChannelID); err != nil {
		Logger.Debug("failed to send typing indicator", zap.Error(err))
		return false
	}

	newPrompt := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: data.Content,
	}

	numNewPromptToken := predictTokens([]openai.ChatCompletionMessage{newPrompt}, false)
	remainingTokens := channelConfig.PromptTokenLimit - 3 - numNewPromptToken

	if remainingTokens < 0 {
		sendErrorMessage(s, data, serverConfig.Language, "token_limit_reached")
		return true
	}

	var prompts []openai.ChatCompletionMessage
	previousMessages, tokens, fetchErr := MessageDatabase.Fetch(data.Author.ID, remainingTokens)
	if fetchErr != nil {
		Logger.Debug("failed to fetch previous messages", zap.Error(fetchErr))
		return true
	}

	for i := len(previousMessages) - 1; i >= 0; i-- {
		prompts = append(prompts, *previousMessages[i])
	}

	prompts = append(prompts, newPrompt)

	Logger.Debug(
		"Prompt",
		zap.String("prompts", fmt.Sprintf("%+v", prompts)),
		zap.Int("numNewPromptToken", numNewPromptToken),
	)

	stream, openAIChatErr := OpenAIClient.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
		MaxTokens: channelConfig.CompletionTokenLimit,
		Model:     Model.ID,
		Messages:  prompts,
		User:      data.Author.ID,
	})

	if openAIChatErr != nil {
		sendErrorMessage(s, data, serverConfig.Language, "error_response")
		Logger.Debug("failed to chat with OpenAI", zap.Error(openAIChatErr))
		return true
	}

	defer stream.Close()

	responseMessage, numResponseMessage, discordResponseErr := sendDiscordResponseWithStream(
		stream, time.Duration(channelConfig.MessageEditInterval)*time.Millisecond, s,
		data.GuildID, data.ChannelID, data.ID, serverConfig.Language, tokens+numNewPromptToken+3)
	if discordResponseErr != nil {
		sendErrorMessage(s, data, serverConfig.Language, "error_response")
		Logger.Debug("failed to send Discord response", zap.Error(discordResponseErr))
		return true
	}

	// Store the bot response in the database
	if err := storeInteraction(
		data.Author.ID,
		&newPrompt, numNewPromptToken,
		responseMessage, numResponseMessage,
	); err != nil {
		Logger.Debug("failed to store interaction", zap.Error(err))
	}

	return true
}

func init() {
	messageHandlers = append(messageHandlers, chatChanel)
}
