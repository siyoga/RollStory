package context

import (
	"context"

	"github.com/siyoga/rollstory/internal/errors"
)

func (s *service) CreateThreadAndSendInstruction(ctx context.Context, userId int64) (string, *errors.Error) {
	threadId, err := s.threadRepository.GetThreadByUser(ctx, userId)
	if err != nil {
		return "", errors.DatabaseError(err)
	}

	if threadId == "" {
		thread, err := s.gptAdapter.CreateThread(ctx)
		if err != nil {
			return "", errors.AdapterError(err)
		}

		if err := s.threadRepository.SaveThreadForUser(ctx, thread.ID, userId); err != nil {
			return "", errors.DatabaseError(err)
		}
	}

	return "Чтобы начать игру нужно:\n1. Задать мир с помощью команды /world\n2. Задать персонажа с помощью команды /character", nil
}
