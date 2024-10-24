package gpt

import (
	"context"

	"github.com/siyoga/rollstory/internal/config"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/logger"

	"github.com/sashabaranov/go-openai"
)

type (
	Adapter interface {
		CreateThread(ctx context.Context) (openai.Thread, error)
		Request(ctx context.Context, threadId string, msg string, respLimit int, respOrder domain.ReturnOrder) (openai.MessagesList, error)
	}

	adapter struct {
		cfg    config.OpenAI
		log    logger.Logger
		client *openai.Client
	}
)

func NewAdapter(
	cfg config.OpenAI,
	log logger.Logger,
) *adapter {
	return &adapter{
		cfg:    cfg,
		log:    log,
		client: openai.NewClient(cfg.Token),
	}
}
