package adapter

import (
	"context"

	"github.com/siyoga/rollstory/internal/domain"

	"github.com/sashabaranov/go-openai"
)

type (
	TelegramAdapter interface {
		SendMessage(chatId int64, text string) error
		Updates(offset int, limit int) ([]domain.Update, error)
	}

	OpenAIAdapter interface {
		CreateThread(ctx context.Context) (openai.Thread, error)
		Request(ctx context.Context, threadId string, msg string, respLimit int, respOrder domain.ReturnOrder) (openai.MessagesList, error)
	}
)
