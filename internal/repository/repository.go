package repository

import (
	"context"

	"github.com/siyoga/rollstory/internal/models"
)

type (
	UserRepository interface {
		SaveUser(ctx context.Context, userId int, user models.User) error
		GetUser(ctx context.Context, userID int) (models.User, error)
	}
)
