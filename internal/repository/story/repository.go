package story

import (
	"github.com/redis/go-redis/v9"
	"github.com/siyoga/rollstory/internal/app/db"
	"github.com/siyoga/rollstory/internal/logger"
)

type repository struct {
	log    logger.Logger
	client *redis.Client
}

func NewThreadRepository(log logger.Logger, client *db.RedisClient) *repository {
	return &repository{
		log:    log.Named("story_repo"),
		client: client.DB,
	}
}
