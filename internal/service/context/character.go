package context

import (
	"context"
	"fmt"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
	"github.com/siyoga/rollstory/internal/models"
)

func (s *service) CreateCharacter(ctx context.Context, userId int, user *domain.UserInfo, characterDesc string) (string, *errors.Error) {
	if _, err := s.gptAdapter.Request(ctx, user.ThreadId, fmt.Sprintf("Описание персонажа:\n\n%s", characterDesc), 1, domain.Asc); err != nil {
		return "", errors.AdapterError(err)
	}

	user.Character = characterDesc
	if err := s.userRepository.SaveUser(ctx, userId, models.User{}.FromDomain(*user)); err != nil {
		return "", errors.DatabaseError(err)
	}

	return "✅Персонаж сохранён!", nil
}

func (s *service) EditCharacter(ctx context.Context, userId int, user *domain.UserInfo, newCharacterDesc string) (string, *errors.Error) {
	if _, err := s.gptAdapter.Request(ctx, user.ThreadId, fmt.Sprintf("Забудь старое описание персонажа. Это новое описание персонажа:\n\n%s", newCharacterDesc), 1, domain.Asc); err != nil {
		return "", errors.AdapterError(err)
	}

	user.Character = newCharacterDesc
	if err := s.userRepository.SaveUser(ctx, userId, models.User{}.FromDomain(*user)); err != nil {
		return "", errors.DatabaseError(err)
	}

	return "✅Персонаж обновлен!", nil
}
