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

	story, e := s.getStoryByUser(ctx, userId)
	if e != nil {
		return "", e
	}

	if story.Character != "" {
		return fmt.Sprintf("Описание персонажа уже установлено:\n%s\n\nУстановите описание мира (/world), чтобы начать игру.", story.Character), nil
	}

	story.Character = characterDesc
	if err := s.storyRepository.SaveSettingsByUser(ctx, userId, story.ToModel()); err != nil {
		return "", errors.DatabaseError(err)
	}

	return resp.Messages[0].Content[0].Text.Value, nil
}
