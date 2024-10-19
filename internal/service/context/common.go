package context

import (
	"context"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (s *service) getStoryByUser(ctx context.Context, userId int64) (story domain.Story, e *errors.Error) {
	m, err := s.storyRepository.GetStoryByUser(ctx, userId)
	if err != nil {
		return domain.Story{}, errors.DatabaseError(err)
	}

	if story.FromModel(m).IsEmpty() {
		if err := s.storyRepository.SaveSettingsByUser(ctx, userId, domain.Story{}.ToModel()); err != nil {
			return domain.Story{}, errors.DatabaseError(err)
		}

		return domain.Story{}, nil
	}

	return domain.Story{}.FromModel(m), nil
}
