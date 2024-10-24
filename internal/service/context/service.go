package context

import (
	"github.com/siyoga/rollstory/internal/adapter/gpt"
	"github.com/siyoga/rollstory/internal/logger"
	"github.com/siyoga/rollstory/internal/repository"
	def "github.com/siyoga/rollstory/internal/service"
)

var _ def.ContextService = (*service)(nil)

type service struct {
	log        logger.Logger
	gptAdapter gpt.Adapter

	storyRepository  repository.StoryRepository
	threadRepository repository.ThreadRepository
}

func NewContextService(
	log logger.Logger,
	gptAdapter gpt.Adapter,

	storyRepository repository.StoryRepository,
	threadRepository repository.ThreadRepository,
) *service {
	return &service{
		log:        log,
		gptAdapter: gptAdapter,

		storyRepository:  storyRepository,
		threadRepository: threadRepository,
	}
}
