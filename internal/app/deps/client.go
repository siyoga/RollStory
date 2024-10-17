package deps

import (
	"github.com/siyoga/rollstory/internal/app/db"
	"go.uber.org/zap"
)

func (d *dependencies) RedisClient() *db.RedisClient {
	if d.redisClient == nil {
		var err error
		msg := "initializing redis client"
		if d.redisClient, err = db.NewRedisClient(d.cfg.Redis.DSN, d.cfg.Redis.CertLoc); err != nil {
			d.log.Zap().Panic(msg, zap.Error(err))
		}

		d.closeCallbacks = append(d.closeCallbacks, func() {
			if err := d.redisClient.Close(); err != nil {
				msg := "stop redis client"
				d.log.Zap().Warn(msg, zap.Error(err))
				return
			}
		})
	}

	return d.redisClient
}
