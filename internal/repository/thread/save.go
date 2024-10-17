package thread

import (
	"context"
	"strconv"

	"github.com/siyoga/rollstory/internal/errors"
)

func (r *repository) SaveThreadForUser(ctx context.Context, threadId string, userId int64) error {
	if err := r.client.Set(ctx, strconv.Itoa(int(userId)), threadId, 0).Err(); err != nil {
		return r.log.DbError(err, errors.ErrRedisSaveRaw, "redis save")
	}

	return nil
}
