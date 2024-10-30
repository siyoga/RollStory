package game

import (
	"context"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
	"github.com/siyoga/rollstory/internal/models"
)

func (s *service) GameMessage(ctx context.Context, userId int, user *domain.UserInfo, message string) (string, *errors.Error) {
	resp, err := s.gptAdapter.Request(ctx, user.ThreadId, message, 1, domain.Asc)
	if err != nil {
		return "", errors.AdapterError(err)
	}

	return resp.Messages[0].Content[0].Text.Value, nil
}

func (s *service) NewGame(ctx context.Context, userId int, user *domain.UserInfo) (string, *errors.Error) {
	if err := s.gptAdapter.DeleteThread(ctx, user.ThreadId); err != nil {
		return "", errors.AdapterError(err)
	}

	thread, err := s.gptAdapter.CreateThread(ctx)
	if err != nil {
		return "", errors.AdapterError(err)
	}

	user.ThreadId = thread.ID
	user.World = ""
	user.Character = ""
	user.IsStarted = false

	if err := s.userRepository.SaveUser(ctx, userId, models.User{}.FromDomain(*user)); err != nil {
		return "", errors.DatabaseError(err)
	}

	return "Текущая игра остановлена и не может быть продолжена. Для начала новой игры задайте игровой мир и персонажа.", nil
}
