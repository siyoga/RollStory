package user

import (
	"context"
	"encoding/json"
	errs "errors"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/siyoga/rollstory/internal/errors"
	"github.com/siyoga/rollstory/internal/models"
)

func (r *repository) GetUser(ctx context.Context, userId int) (models.User, error) {
	data, err := r.client.Get(ctx, strconv.Itoa(userId)).Result()
	if err != nil {
		if errs.Is(err, redis.Nil) {
			return models.User{}, nil
		}

		return models.User{}, r.log.DbError(err, errors.ErrRedisGetRaw, "redis get")
	}

	var story models.User
	if err := json.Unmarshal([]byte(data), &story); err != nil {
		return models.User{}, r.log.DbError(err, errors.ErrRedisGetRaw, "redis get")
	}

	return story, nil
}
