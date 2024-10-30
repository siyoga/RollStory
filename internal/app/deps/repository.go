package deps

import (
	"github.com/siyoga/rollstory/internal/repository"
	"github.com/siyoga/rollstory/internal/repository/user"
)

func (d *dependencies) UserRepository() repository.UserRepository {
	if d.userRepo == nil {
		d.userRepo = user.NewUserRepository(d.log, d.RedisClient())
	}

	return d.userRepo
}
