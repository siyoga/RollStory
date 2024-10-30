package context

import (
	"context"

	"github.com/siyoga/rollstory/internal/errors"
	"github.com/siyoga/rollstory/internal/models"
)

func (s *service) CreateThreadAndSendInstruction(ctx context.Context, userId int) (string, *errors.Error) {
	user, err := s.userRepository.GetUser(ctx, userId)
	if err != nil {
		return "", errors.DatabaseError(err)
	}

	if user.ToDomain().IsEmpty() {
		thread, err := s.gptAdapter.CreateThread(ctx)
		if err != nil {
			return "", errors.AdapterError(err)
		}

		if err := s.userRepository.SaveUser(ctx, userId, models.User{ThreadId: thread.ID}); err != nil {
			return "", errors.DatabaseError(err)
		}
	}

	return "Чтобы начать игру нужно:\n1. Задать мир с помощью команды /world\n2. Задать персонажа с помощью команды /character", nil
}
