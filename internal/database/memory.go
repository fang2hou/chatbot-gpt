package database

import (
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// messageData is a struct for storing a message and its timestamp
type messageData struct {
	Message   openai.ChatCompletionMessage
	Token     int
	Timestamp time.Time
}

// MemoryChatDatabase is a simple in-memory database for storing chat messages.
type MemoryChatDatabase struct {
	data map[string][]messageData
}

// Fetch fetches the messages that not exceed the token limit.
func (m *MemoryChatDatabase) Fetch(
	userID string,
	maxToken int,
) ([]*openai.ChatCompletionMessage, int, error) {
	tokens := 0
	var messages []*openai.ChatCompletionMessage

	for _, data := range m.data[userID] {
		thisData := data

		if tokens+thisData.Token < maxToken {
			tokens += thisData.Token
			messages = append(messages, &thisData.Message)
		} else {
			break
		}
	}

	return messages, tokens, nil
}

// Store stores the message in the database
func (m *MemoryChatDatabase) Store(
	userID string,
	newMessage *openai.ChatCompletionMessage,
	numToken int,
) error {
	// Insert the new message at the beginning of the slice
	newMessages := []messageData{
		{
			Message:   *newMessage,
			Token:     numToken,
			Timestamp: time.Now(),
		},
	}

	if m.data[userID] != nil {
		newMessages = append(newMessages, m.data[userID]...)
	}

	m.data[userID] = newMessages

	return nil
}

// Optimize deletes old messages from the database.
func (m *MemoryChatDatabase) Optimize(userID string, tokenLimit int) {
	tokens := 0

	var newMessages []messageData

	for _, message := range m.data[userID] {
		thisMessage := message

		if tokens+thisMessage.Token < tokenLimit {
			tokens += thisMessage.Token
			newMessages = append(newMessages, thisMessage)
		} else {
			break
		}
	}

	m.data[userID] = newMessages
}

// Clear clears the user's message history
func (m *MemoryChatDatabase) Clear(userID string) {
	m.data[userID] = nil
}

// NewMemoryChatDatabase creates a new MemoryChatDatabase
func NewMemoryChatDatabase() ChatDatabase {
	return &MemoryChatDatabase{
		data: make(map[string][]messageData),
	}
}
