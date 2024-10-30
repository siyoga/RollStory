package game

import (
	"github.com/siyoga/rollstory/internal/adapter/gpt"
	"github.com/siyoga/rollstory/internal/logger"
	"github.com/siyoga/rollstory/internal/repository"
	def "github.com/siyoga/rollstory/internal/service"
)

var _ def.GameService = (*service)(nil)

type service struct {
	log            logger.Logger
	gptAdapter     gpt.Adapter
	userRepository repository.UserRepository
}

func NewGameService(
	log logger.Logger,
	gptAdapter gpt.Adapter,
	userRepository repository.UserRepository,
) *service {
	return &service{
		log:            log,
		gptAdapter:     gptAdapter,
		userRepository: userRepository,
	}
}
