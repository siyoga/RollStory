package context

import (
	"github.com/siyoga/rollstory/internal/adapter"
	"github.com/siyoga/rollstory/internal/logger"
	"github.com/siyoga/rollstory/internal/repository"
	def "github.com/siyoga/rollstory/internal/service"
)

var _ def.ContextService = (*service)(nil)

type service struct {
	log              logger.Logger
	gptAdapter       adapter.OpenAIAdapter
	threadRepository repository.ThreadRepository
}

func NewContextService(
	log logger.Logger,
	gptAdapter adapter.OpenAIAdapter,
	threadRepository repository.ThreadRepository,
) *service {
	return &service{
		log:              log,
		gptAdapter:       gptAdapter,
		threadRepository: threadRepository,
	}
}
