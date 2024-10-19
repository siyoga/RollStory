package router

import (
	"context"

	"github.com/siyoga/rollstory/internal/config"
	"github.com/siyoga/rollstory/internal/domain"
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
) Handler {
	return &contextHandler{
		timeouts: timeouts,
		service:  ctxService,
	}
}

func (c *contextHandler) FillHandlers(r Router) {
	r.Handle("start", c.start)
	r.Handle("character", c.character)
	r.Handle("world", c.world)
	r.Handle("begin")
}

func (c *contextHandler) start(ctx context.Context, userId int64, msg *domain.Message) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	res, e := c.service.CreateThreadAndSendInstruction(ctx, userId)
	if e != nil {
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.MessageResult{
			Message: res,
			ChatId:  msg.Chat.ID,
		},
		201,
		&userId,
	)
}

func (c *contextHandler) character(ctx context.Context, userId int64, msg *domain.Message) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	if msg.IsCommand() {
		return newSuccessResponse(
			domain.MessageResult{
				Message: "Отправьте описание вашего персонажа, по пунктам: \n" +
					"1. Внешность\n" +
					"2. Характер\n" +
					"3. Краткая предыстория, его цели, мотивация\n" +
					"4. Имя\n" +
					"5. Ключевые персонажи/группы людей (враги, союзники, нейтральные)",
				ChatId: msg.Chat.ID,
			},
			200,
			nil,
		)
	}

	res, e := c.service.CreateCharacter(ctx, userId, msg.Text)
	if e != nil {
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.MessageResult{
			Message: res,
			ChatId:  msg.Chat.ID,
		},
		200,
		&userId,
	)
}

func (c *contextHandler) world(ctx context.Context, userId int64, msg *domain.Message) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	if msg.IsCommand() {
		return newSuccessResponse(
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
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.MessageResult{
			Message: res,
			ChatId:  msg.Chat.ID,
		},
		200,
		&userId,
	)
}

func (c *contextHandler) begin(ctx context.Context, userId int64, msg *domain.Message) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	res, e := c.service.CreateWorld(ctx, userId, msg.Text)
	if e != nil {
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.MessageResult{
			Message: res,
			ChatId:  msg.Chat.ID,
		},
		200,
		&userId,
	)
}
