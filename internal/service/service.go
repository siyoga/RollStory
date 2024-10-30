package service

import (
	"context"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

type (
	ContextService interface {
		CreateThreadAndSendInstruction(ctx context.Context, userId int) (string, *errors.Error)
		CreateWorld(ctx context.Context, userId int, user *domain.UserInfo, worldDesc string) (string, *errors.Error)
		CreateCharacter(ctx context.Context, userId int, user *domain.UserInfo, characterDesc string) (string, *errors.Error)

		EditCharacter(ctx context.Context, userId int, user *domain.UserInfo, newCharacterDesc string) (string, *errors.Error)
		EditWorld(ctx context.Context, userId int, user *domain.UserInfo, newWorldDesc string) (string, *errors.Error)

		GetUser(ctx context.Context, userId int) (domain.UserInfo, *errors.Error)
		BeginStory(ctx context.Context, userId int, user *domain.UserInfo) (string, *errors.Error)
	}

	GameService interface {
		GameMessage(ctx context.Context, userId int, user *domain.UserInfo, message string) (string, *errors.Error)
		NewGame(ctx context.Context, userId int, user *domain.UserInfo) (string, *errors.Error)
	}
)
