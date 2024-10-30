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

	userRepository repository.UserRepository
}

func NewContextService(
	log logger.Logger,
	gptAdapter gpt.Adapter,

	userRepository repository.UserRepository,
) *service {
	return &service{
		log:        log,
		gptAdapter: gptAdapter,

		userRepository: userRepository,
	}
}
