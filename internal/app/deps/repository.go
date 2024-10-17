package deps

import (
	"github.com/siyoga/rollstory/internal/repository"
	"github.com/siyoga/rollstory/internal/repository/thread"
)

func (d *dependencies) ThreadRepository() repository.ThreadRepository {
	if d.threadRepository == nil {
		d.threadRepository = thread.NewThreadRepository(d.log, d.RedisClient())
	}

	return d.threadRepository
}
