package story

import (
	"context"
	"encoding/json"
	errs "errors"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/siyoga/rollstory/internal/errors"
	"github.com/siyoga/rollstory/internal/models"
)

func (r *repository) GetStoryByUser(ctx context.Context, userId int64) (models.Story, error) {
	data, err := r.client.Get(ctx, strconv.FormatInt(userId, 10)).Bytes()
	if err != nil {
		if errs.Is(err, redis.Nil) {
			return models.Story{}, nil
		}

		return models.Story{}, r.log.DbError(err, errors.ErrRedisGetRaw, "redis get")
	}

	var story models.Story
	if err := json.Unmarshal(data, &story); err != nil {
		return models.Story{}, r.log.DbError(err, errors.ErrRedisGetRaw, "redis get")
	}

	return story, nil
}
