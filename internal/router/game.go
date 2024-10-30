package router

import (
	"context"

	"github.com/siyoga/rollstory/internal/config"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/service"
)

type (
	gameHandler struct {
		timeouts    config.Timeouts
		gameService service.GameService
		ctxService  service.ContextService
	}
)

func NewGameHandler(
	timeouts config.Timeouts,
	gameService service.GameService,
	ctxService service.ContextService,
) Handler {
	return &gameHandler{
		timeouts:    timeouts,
		gameService: gameService,
		ctxService:  ctxService,
	}
}

func (g *gameHandler) FillHandlers(r Router) {
	r.DefaultRoute(g.message)
	r.Route(NewGameSlash, NewGamePlain).Handle(g.newGame)
}

func (g *gameHandler) message(ctx context.Context, userId int, req domain.Request) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, g.timeouts.RequestTimeout)
	defer cancel()

	user, e := g.ctxService.GetUser(ctx, userId)
	if e != nil {
		return newErrResponse(e, userId)
	}

	answer, e := g.gameService.GameMessage(ctx, userId, &user, req.Data)
	if e != nil {
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.Request{
			Data:   answer,
			ChatId: req.ChatId,
			Markup: keyboardMiddleware(user),
		},
		&userId,
	)
}

func (g *gameHandler) newGame(ctx context.Context, userId int, req domain.Request) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, g.timeouts.RequestTimeout)
	defer cancel()

	user, e := g.ctxService.GetUser(ctx, userId)
	if e != nil {
		return newErrResponse(e, userId)
	}

	answer, e := g.gameService.NewGame(ctx, userId, &user)
	if e != nil {
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.Request{
			Data:   answer,
			ChatId: req.ChatId,
			Markup: keyboardMiddleware(user),
		},
		&userId,
	)
}
