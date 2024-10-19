package deps

import (
	"github.com/siyoga/rollstory/internal/app/db"
	"go.uber.org/zap"
)

func (d *dependencies) RedisThreadClient() *db.RedisClient {
	if d.redisThreadClient == nil {
		var err error
		msg := "initializing redis client"
		if d.redisThreadClient, err = db.NewRedisClient(d.cfg.Redis.ThreadDSN, d.cfg.Redis.CertLoc); err != nil {
			d.log.Zap().Panic(msg, zap.Error(err))
		}

		d.closeCallbacks = append(d.closeCallbacks, func() {
			if err := d.redisThreadClient.Close(); err != nil {
				msg := "stop redis client"
				d.log.Zap().Warn(msg, zap.Error(err))
				return
			}
		})
	}

	return d.redisThreadClient
}

func (d *dependencies) RedisStoryClient() *db.RedisClient {
	if d.redisThreadClient == nil {
		var err error
		msg := "initializing redis client"
		if d.redisThreadClient, err = db.NewRedisClient(d.cfg.Redis.StoryDSN, d.cfg.Redis.CertLoc); err != nil {
			d.log.Zap().Panic(msg, zap.Error(err))
		}

		d.closeCallbacks = append(d.closeCallbacks, func() {
			if err := d.redisThreadClient.Close(); err != nil {
				msg := "stop redis client"
				d.log.Zap().Warn(msg, zap.Error(err))
				return
			}
		})
	}

	return d.redisThreadClient
}
