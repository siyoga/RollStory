package thread

import (
	"context"
	internal_errors "errors"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/siyoga/rollstory/internal/errors"
)

func (r *repository) GetThreadByUser(ctx context.Context, userId int64) (string, error) {
	threadId, err := r.client.Get(ctx, strconv.Itoa(int(userId))).Result()
	if err != nil {
		if internal_errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", r.log.DbError(err, errors.ErrRedisGetRaw, "redis get")
	}

	return threadId, nil
}
