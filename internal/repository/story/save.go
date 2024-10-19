package story

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/siyoga/rollstory/internal/errors"
	"github.com/siyoga/rollstory/internal/models"
)

func (r *repository) SaveSettingsByUser(ctx context.Context, userId int64, story models.Story) error {
	obj, err := json.Marshal(story)
	if err != nil {
		return r.log.DbError(err, errors.ErrRedisSaveRaw, "redis save")
	}

	if err := r.client.Set(ctx, strconv.FormatInt(userId, 10), obj, 0).Err(); err != nil {
		return r.log.DbError(err, errors.ErrRedisSaveRaw, "redis save")
	}

	return nil
}
