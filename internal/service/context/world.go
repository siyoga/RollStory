package context

import (
	"context"
	"fmt"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
	"github.com/siyoga/rollstory/internal/models"
)

func (s *service) CreateWorld(ctx context.Context, userId int, user *domain.UserInfo, worldDesc string) (string, *errors.Error) {
	if _, err := s.gptAdapter.Request(ctx, user.ThreadId, fmt.Sprintf("Описание мира:\n%s", worldDesc), 1, domain.Asc); err != nil {
		return "", errors.AdapterError(err)
	}

	user.World = worldDesc
	if err := s.userRepository.SaveUser(ctx, userId, models.User{}.FromDomain(*user)); err != nil {
		return "", errors.DatabaseError(err)
	}

	return "✅Контекст игрового мира сохранён!", nil
}

func (s *service) EditWorld(ctx context.Context, userId int, user *domain.UserInfo, newWorldDesc string) (string, *errors.Error) {
	if _, err := s.gptAdapter.Request(ctx, user.ThreadId, fmt.Sprintf("Забудь старое описание мира. Это новое описание мира:\n%s", newWorldDesc), 1, domain.Asc); err != nil {
		return "", errors.AdapterError(err)
	}

	user.World = newWorldDesc
	if err := s.userRepository.SaveUser(ctx, userId, models.User{}.FromDomain(*user)); err != nil {
		return "", errors.DatabaseError(err)
	}

	return "✅Контекст игрового мира сохранён!", nil
}
