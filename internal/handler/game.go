package handler

import (
	"context"

	"github.com/siyoga/rollstory/internal/config"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/handler/router"
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
	router router.Router,
) Handler {
	g := gameHandler{
		timeouts: timeouts,
		service:  gameService,
	}

	g.fillHandlers(router)

	return &g
}

func (g *gameHandler) fillHandlers(r router.Router) {
	r.DefaultHandle(g.messageHandler)
}

func (g *gameHandler) messageHandler(ctx context.Context, userId int64, msg *domain.Message) router.Response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, g.timeouts.RequestTimeout)
	defer cancel()

	answer, e := g.service.GameMessage(ctx, userId, msg.Text)
	if e != nil {
		return router.NewErrResponse(e, userId)
	}

	return router.NewSuccessResponse(
		domain.MessageResult{
			Message: answer,
			ChatId:  msg.Chat.ID,
		},
		200,
		&userId,
	)
}
