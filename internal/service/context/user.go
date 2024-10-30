package context

import (
	"context"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (s *service) GetUser(ctx context.Context, userId int) (domain.UserInfo, *errors.Error) {
	user, err := s.userRepository.GetUser(ctx, userId)
	if err != nil {
		return domain.UserInfo{}, errors.DatabaseError(err)
	}

	return user.ToDomain(), nil
}
