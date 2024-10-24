package router

import (
	"context"

	"github.com/siyoga/rollstory/internal/config"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/service"
)

type (
	gameHandler struct {
		timeouts config.Timeouts
		service  service.GameService
	}
)

func NewGameHandler(
	timeouts config.Timeouts,
	gameService service.GameService,
) Handler {
	return &gameHandler{
		timeouts: timeouts,
		service:  gameService,
	}
}

func (g *gameHandler) FillHandlers(r Router) {
	r.DefaultHandle(g.messageHandler)
}

func (g *gameHandler) messageHandler(ctx context.Context, userId int64, msg *domain.Message) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, g.timeouts.RequestTimeout)
	defer cancel()

	answer, e := g.service.GameMessage(ctx, userId, msg.Text)
	if e != nil {
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.MessageResult{
			Message: answer,
			ChatId:  msg.Chat.ID,
		},
		200,
		&userId,
	)
}
