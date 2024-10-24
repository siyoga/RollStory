package context

import (
	"context"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (s *service) BeginStory(ctx context.Context, userId int64) (string, *errors.Error) {
	settings, err := s.storyRepository.GetStoryByUser(ctx, userId)
	if err != nil {
		return "", errors.DatabaseError(err)
	}

	if settings.Character == "" {
		return "Не установлено описание игрового персонажа (/character). Установите описание персонажа, и начните игру еще раз.", nil
	}

	if settings.World == "" {
		return "Не установлено описание игрового мира (/world). Установите описание мира, и начните игру еще раз.", nil
	}

	threadId, err := s.threadRepository.GetThreadByUser(ctx, userId)
	if err != nil {
		return "", errors.DatabaseError(err)
	}

	resp, err := s.gptAdapter.Request(ctx, threadId, "Начинаем игру", 1, domain.Asc)
	if err != nil {
		return "", errors.AdapterError(err)
	}

	return resp.Messages[0].Content[0].Text.Value, nil
}
