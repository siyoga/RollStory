package repository

import "context"

type (
	ThreadRepository interface {
		SaveThreadForUser(ctx context.Context, threadId string, userId int64) error
		GetThreadByUser(ctx context.Context, userID int64) (string, error)
	}
)
