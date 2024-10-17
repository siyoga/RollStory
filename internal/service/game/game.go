package game

import (
	"context"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (s *service) GameMessage(ctx context.Context, userId int64, message string) (string, *errors.Error) {
	threadId, err := s.threadRepository.GetThreadByUser(ctx, userId)
	if err != nil {
		return "", errors.DatabaseError(err)
	}

	resp, err := s.gptAdapter.Request(ctx, threadId, message, 1, domain.Asc)
	if err != nil {
		return "", errors.AdapterError(err)
	}

	return resp.Messages[0].Content[0].Text.Value, nil
}
