package repository

import (
	"context"

	"github.com/siyoga/rollstory/internal/models"
)

type (
	ThreadRepository interface {
		SaveThreadForUser(ctx context.Context, threadId string, userId int64) error
		GetThreadByUser(ctx context.Context, userID int64) (string, error)
	}

	StoryRepository interface {
		SaveSettingsByUser(ctx context.Context, userId int64, story models.Story) error
		GetStoryByUser(ctx context.Context, userId int64) (models.Story, error)
	}
)
