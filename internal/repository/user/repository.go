package user

import (
	"github.com/redis/go-redis/v9"
	"github.com/siyoga/rollstory/internal/app/db"
	"github.com/siyoga/rollstory/internal/logger"
	def "github.com/siyoga/rollstory/internal/repository"
)

var _ def.UserRepository = (*repository)(nil)

type repository struct {
	log    logger.Logger
	client *redis.Client
}

func NewUserRepository(log logger.Logger, client *db.RedisClient) *repository {
	return &repository{
		log:    log.Named("user_repo"),
		client: client.DB,
	}
}
