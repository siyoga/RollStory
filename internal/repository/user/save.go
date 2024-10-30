package user

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/siyoga/rollstory/internal/errors"
	"github.com/siyoga/rollstory/internal/models"
)

func (r *repository) SaveUser(ctx context.Context, userId int, user models.User) error {
	obj, err := json.Marshal(user)
	if err != nil {
		return r.log.DbError(err, errors.ErrRedisSaveRaw, "redis save")
	}

	if err := r.client.Set(ctx, strconv.Itoa(userId), obj, 0).Err(); err != nil {
		return r.log.DbError(err, errors.ErrRedisSaveRaw, "redis save")
	}

	return nil
}
