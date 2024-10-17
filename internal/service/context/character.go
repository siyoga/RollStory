package context

import (
	"context"
	"fmt"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (s *service) CreateCharacter(ctx context.Context, userId int64, characterDesc string) (string, *errors.Error) {
	threadId, err := s.threadRepository.GetThreadByUser(ctx, userId)
	if err != nil {
		return "", errors.DatabaseError(err)
	}

	characterDesc = fmt.Sprintf("Описание персонажа:\n%s", characterDesc)

	resp, err := s.gptAdapter.Request(ctx, threadId, characterDesc, 1, domain.Asc)
	if err != nil {
		return "", errors.AdapterError(err)
	}

	return resp.Messages[0].Content[0].Text.Value, nil
}
