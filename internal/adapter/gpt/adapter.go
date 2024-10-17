package gpt

import (
	def "github.com/siyoga/rollstory/internal/adapter"
	"github.com/siyoga/rollstory/internal/config"
	"github.com/siyoga/rollstory/internal/logger"

	"github.com/sashabaranov/go-openai"
)

var _ def.OpenAIAdapter = (*adapter)(nil)

type adapter struct {
	cfg    config.OpenAI
	log    logger.Logger
	client *openai.Client
}

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
