package context

import (
	"context"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
	"github.com/siyoga/rollstory/internal/models"
)

func (s *service) BeginStory(ctx context.Context, userId int, user *domain.UserInfo) (string, *errors.Error) {
	if user.Character == "" {
		return "Не установлено описание игрового персонажа (/character). Установите описание персонажа, и начните игру еще раз.", nil
	}

	if user.World == "" {
		return "Не установлено описание игрового мира (/world). Установите описание мира, и начните игру еще раз.", nil
	}

	resp, err := s.gptAdapter.Request(ctx, user.ThreadId, "Начинаем игру", 1, domain.Asc)
	if err != nil {
		return "", errors.AdapterError(err)
	}

	user.IsStarted = true
	if err := s.userRepository.SaveUser(ctx, userId, models.User{}.FromDomain(*user)); err != nil {
		return "", errors.DatabaseError(err)
	}

	return resp.Messages[0].Content[0].Text.Value, nil
}
