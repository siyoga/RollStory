package service

import (
	"context"

	"github.com/siyoga/rollstory/internal/errors"
)

type (
	ContextService interface {
		CreateThreadAndSendInstruction(ctx context.Context, userId int64) (string, *errors.Error)
		CreateWorld(ctx context.Context, userId int64, worldDesc string) (string, *errors.Error)
		CreateCharacter(ctx context.Context, userId int64, characterDesc string) (string, *errors.Error)
	}

	GameService interface {
		GameMessage(ctx context.Context, userId int64, message string) (string, *errors.Error)
	}
)
