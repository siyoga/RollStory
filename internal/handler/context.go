package handler

import (
	"context"

	"github.com/siyoga/rollstory/internal/config"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/handler/router"
	"github.com/siyoga/rollstory/internal/service"
)

type (
	contextHandler struct {
		timeouts config.Timeouts
		service  service.ContextService
	}
)

func NewContextHandler(
	timeouts config.Timeouts,
	ctxService service.ContextService,
	router router.Router,
) Handler {
	c := contextHandler{
		service: ctxService,
	}

	c.fillHandlers(router)

	return &c
}

func (c *contextHandler) fillHandlers(r router.Router) {
	r.Handle("start", c.start)
	r.Handle("character", c.character)
	r.Handle("world", c.world)
}

func (c *contextHandler) start(ctx context.Context, userId int64, msg *domain.Message) router.Response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	res, e := c.service.CreateThreadAndSendInstruction(ctx, userId)
	if e != nil {
		return router.NewErrResponse(e, userId)
	}

	return router.NewSuccessResponse(
		domain.MessageResult{
			Message: res,
			ChatId:  msg.Chat.ID,
		},
		201,
		&userId,
	)
}

func (c *contextHandler) character(ctx context.Context, userId int64, msg *domain.Message) router.Response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	if msg.IsCommand() {
		return router.NewSuccessResponse(
			domain.MessageResult{
				Message: "Отправьте описание вашего персонажа, желательно чтобы в нем было: \n" +
					"1. Внешность\n" +
					"2. Характер\n" +
					"3. Краткая предыстория, его цели, мотивация",
				ChatId: msg.Chat.ID,
			},
			200,
			nil,
		)
	}

	res, e := c.service.CreateCharacter(ctx, userId, msg.Text)
	if e != nil {
		return router.NewErrResponse(e, userId)
	}

	return router.NewSuccessResponse(
		domain.MessageResult{
			Message: res,
			ChatId:  msg.Chat.ID,
		},
		200,
		&userId,
	)
}

func (c *contextHandler) world(ctx context.Context, userId int64, msg *domain.Message) router.Response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	if msg.IsCommand() {
		return router.NewSuccessResponse(
			domain.MessageResult{
				Message: "Отправьте жанр вселенной в которой хотите играть," +
					" чем подробнее будет описаниее, тем лучше будет ваш игровой опыт",
				ChatId: msg.Chat.ID,
			},
			200,
			nil,
		)
	}

	res, e := c.service.CreateWorld(ctx, userId, msg.Text)
	if e != nil {
		return router.NewErrResponse(e, userId)
	}

	return router.NewSuccessResponse(
		domain.MessageResult{
			Message: res,
			ChatId:  msg.Chat.ID,
		},
		200,
		&userId,
	)
}
