package main

import (
	"context"
	"flag"

	"github.com/bwmarrin/discordgo"
	tiktoken "github.com/pkoukk/tiktoken-go"
	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/zap"

	"chatbot-gpt/internal/config"
	"chatbot-gpt/internal/cost"
	"chatbot-gpt/internal/database"
	"chatbot-gpt/internal/locale"
)

// ChannelConfig is the configuration for a channel.
type ChannelConfig struct {
	MessageEditInterval  int
	PromptTokenLimit     int
	CompletionTokenLimit int
}

// ServerConfig is the configuration for a server.
type ServerConfig struct {
	Language     locale.Language
	ChatChannels map[string]ChannelConfig
	Commands     struct {
		ClearContext struct {
			Enable  bool
			Aliases []string
		}
	}
}

const (
	// ConfigPrefix is the prefix used for handle environment variables.
	configPrefix = "CHATBOT_GPT"
)

var (
	// Logger is the logger used by the bot.
	Logger *zap.Logger

	// OpenAIClient is the OpenAI client used by the bot.
	OpenAIClient *openai.Client

	// DiscordClient is the Discord client used by the bot.
	DiscordClient *discordgo.Session

	// Model is the OpenAI model used by the bot.
	Model *openai.Model

	// TokenPredictionModel is the model id that used by the bot for token prediction.
	TokenPredictionModel *tiktoken.Tiktoken

	// ServerConfigMap is the map of server configurations.
	ServerConfigMap map[string]ServerConfig

	// Localizer is the localizer used to translate messages.
	Localizer locale.Localizer

	// MessageDatabase is the database used to store messages.
	MessageDatabase database.ChatDatabase

	// CostCalculator is the calculator used to calculate the cost of a message.
	CostCalculator *cost.Calculator
)

// initLogger initializes the logger.
func initLogger(isProduction bool) {
	if isProduction {
		if l, err := zap.NewProduction(); err != nil {
			panic(err)
		} else {
			Logger = l
		}
	} else {
		if l, err := zap.NewDevelopment(); err != nil {
			panic(err)
		} else {
			Logger = l
		}
	}
}

// initOpenAIClient initializes the OpenAI client.
func initOpenAIClient(cfg config.OpenAI) {
	if tkm, err := tiktoken.EncodingForModel(cfg.TokenPredictionModelID); err != nil {
		Logger.Panic("failed to initialize token prediction model", zap.Error(err))
	} else {
		TokenPredictionModel = tkm
	}

	OpenAIClient = openai.NewClient(cfg.Token)

	if result, err := OpenAIClient.ListModels(context.Background()); err != nil {
		Logger.Panic("failed to initialize OpenAI client", zap.Error(err))
	} else {
		for _, model := range result.Models {
			thisModel := model
			if model.ID == cfg.ModelID {
				Model = &thisModel
				break
			}
		}

		if Model == nil {
			Logger.Panic("invalid model ID or you have not access to it.", zap.String("modelID", cfg.ModelID))
		}
	}
}

// initDiscordClient initializes the Discord client.
func initDiscordClient(cfg config.Discord) {
	if s, err := discordgo.New("Bot " + cfg.Token); err != nil {
		Logger.Panic("failed to create discord session", zap.Error(err))
	} else {
		DiscordClient = s
	}
}

// initLocalizer initializes the localizer.
func initLocalizer(cfg config.Discord) {
	Localizer = locale.NewLocalizer()
	for key, translations := range cfg.Locales {
		for langCode, translation := range translations {
			if lang, err := locale.ToLanguage(langCode); err != nil {
				Logger.Warn("invalid language code", zap.String("code", langCode))
			} else {
				Localizer.Update(key, lang, translation)
			}
		}
	}
}

// initServerConfigMap initializes the server configuration map.
func initServerConfigMap(cfg config.Discord) {
	ServerConfigMap = make(map[string]ServerConfig)

	for _, serverConfig := range cfg.Servers {
		chatChannels := make(map[string]ChannelConfig)

		for _, channelConfig := range serverConfig.ChatChannels {
			chatChannels[channelConfig.ID] = ChannelConfig{
				MessageEditInterval:  channelConfig.MessageEditInterval,
				PromptTokenLimit:     channelConfig.PromptTokenLimit,
				CompletionTokenLimit: channelConfig.CompletionTokenLimit,
			}
		}

		language, langParseErr := locale.ToLanguage(serverConfig.Language)

		if langParseErr != nil {
			Logger.Panic("invalid language code", zap.String("code", serverConfig.Language))
		}

		ServerConfigMap[serverConfig.ID] = ServerConfig{
			Language:     language,
			ChatChannels: chatChannels,
			Commands: struct {
				ClearContext struct {
					Enable  bool
					Aliases []string
				}
			}(serverConfig.Commands),
		}
	}
}

// initMessageDatabase initializes the message database.
func initMessageDatabase() {
	MessageDatabase = database.NewMemoryChatDatabase()
}

// initCostCalculator initializes the cost calculator.
func initCostCalculator() {
	CostCalculator = cost.NewCalculator(Model)
}

func init() {
	path := flag.String("config", "config.json", "Path to the cfg file")
	flag.Parse()

	userConfig, err := config.Init(&struct {
		Discord config.Discord
		OpenAI  config.OpenAI
	}{}, configPrefix, *path)
	if err != nil {
		panic(err)
	}

	initLogger(userConfig.Discord.Production)
	initMessageDatabase()
	initOpenAIClient(userConfig.OpenAI)
	initCostCalculator()
	initDiscordClient(userConfig.Discord)
	initLocalizer(userConfig.Discord)
	initServerConfigMap(userConfig.Discord)
}
