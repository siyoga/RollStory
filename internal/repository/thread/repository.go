package thread

import (
	"github.com/redis/go-redis/v9"
	"github.com/siyoga/rollstory/internal/app/db"
	"github.com/siyoga/rollstory/internal/logger"
	def "github.com/siyoga/rollstory/internal/repository"
)

var _ def.ThreadRepository = (*repository)(nil)

type repository struct {
	log    logger.Logger
	client *redis.Client
}

func NewThreadRepository(log logger.Logger, client *db.RedisClient) *repository {
	return &repository{
		log:    log.Named("thread_repo"),
		client: client.DB,
	}
}
