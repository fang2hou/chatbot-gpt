package database

import (
	"github.com/sashabaranov/go-openai"
)

type ChatDatabase interface {
	Fetch(userID string, maxToken int) (messages []*openai.ChatCompletionMessage, numToken int, err error)
	Store(userID string, newMessage *openai.ChatCompletionMessage, numToken int) error
	Optimize(userID string, tokenLimit int)
	Clear(userID string)
}
