package router

import (
	"context"
	"fmt"

	"github.com/siyoga/rollstory/internal/config"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/service"
)

// "context" cause store info about upcoming game
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
	r.Route(StartSlash).Handle(c.start)

	r.Route(CharacterSlash, CharacterPlain).Handle(c.character)
	r.Route(EditCharacterSlash, EditCharacterPlain).Handle(c.editCharacter).LinkTo(CharacterSlash, CharacterPlain)

	r.Route(WorldSlash, WorldPlain).Handle(c.world)
	r.Route(EditWorldSlash, EditWorldPlain).Handle(c.editWorld).LinkTo(WorldSlash, WorldPlain)

	r.Route(BeginSlash, BeginPlain).Handle(c.begin)
}

func (c *contextHandler) defaultInlineMarkup() domain.InlineMarkup {
	return domain.InlineMarkup{
		Keyboard: map[int][]domain.Button{
			1: []domain.Button{
				{
					Text: "Отменить",
					Data: "/cancel",
				},
			},
		},
	}
}

func (c *contextHandler) start(ctx context.Context, userId int, req domain.Request) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	user, e := c.service.GetUser(ctx, userId)
	if e != nil {
		return newErrResponse(e, userId)
	}

	res, e := c.service.CreateThreadAndSendInstruction(ctx, userId)
	if e != nil {
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.Request{
			Data:   res,
			ChatId: req.ChatId,
			Markup: keyboardMiddleware(user),
		},
		&userId,
	)
}

func (c *contextHandler) character(ctx context.Context, userId int, req domain.Request) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	user, e := c.service.GetUser(ctx, userId)
	if e != nil {
		return newErrResponse(e, userId)
	}

	defaultMarkup := c.defaultInlineMarkup()

	if user.Character != "" {
		defaultMarkup.AddRow([]domain.Button{{Text: string(EditCharacterPlain), Data: string(EditCharacterSlash)}})

		return newSuccessResponse(
			domain.Request{
				Data: fmt.Sprintf(
					"Игровой персонаж уже задан:\n\n"+
						"%s", user.Character),
				ChatId:  req.ChatId,
				ReplyTo: &req.MessageId,
				Markup:  defaultMarkup,
			},
			nil,
		)
	}

	if req.Command != nil {
		return newSuccessResponse(
			domain.Request{
				Data: "Отправьте описание вашего персонажа, по пунктам: \n" +
					"1. Внешность\n" +
					"2. Характер\n" +
					"3. Краткая предыстория, его цели, мотивация\n" +
					"4. Имя\n" +
					"5. Ключевые персонажи/группы людей (враги, союзники, нейтральные)",
				ChatId:  req.ChatId,
				Markup:  defaultMarkup,
				ReplyTo: &req.MessageId,
			},
			nil,
		)
	}

	res, e := c.service.CreateCharacter(ctx, userId, &user, req.Data)
	if e != nil {
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.Request{
			Data:   res,
			ChatId: req.ChatId,
			Markup: keyboardMiddleware(user),
		},

		&userId,
	)
}

func (c *contextHandler) world(ctx context.Context, userId int, req domain.Request) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	user, e := c.service.GetUser(ctx, userId)
	if e != nil {
		return newErrResponse(e, userId)
	}

	defaultMarkup := c.defaultInlineMarkup()

	if user.World != "" {
		defaultMarkup.AddRow([]domain.Button{{Text: string(EditWorldPlain), Data: string(EditWorldSlash)}})

		return newSuccessResponse(
			domain.Request{
				Data: fmt.Sprintf(
					"Игровой мир уже задан:\n\n"+
						"%s", user.World),
				ChatId:  req.ChatId,
				Markup:  defaultMarkup,
				ReplyTo: &req.MessageId,
			},
			nil,
		)
	}

	if req.Command != nil {
		keyboard := make(map[int][]domain.Button)
		keyboard[1] = []domain.Button{
			{
				Text: "Отменить",
				Data: string(domain.Cancel),
			},
		}

		return newSuccessResponse(
			domain.Request{
				Data: "Отправьте жанр вселенной в которой хотите играть," +
					" чем подробнее будет описаниее, тем лучше будет ваш игровой опыт",
				ChatId: req.ChatId,
				Markup: domain.InlineMarkup{
					Keyboard: keyboard,
				},
				ReplyTo: &req.MessageId,
			},
			nil,
		)
	}

	res, e := c.service.CreateWorld(ctx, userId, &user, req.Data)
	if e != nil {
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.Request{
			Data:   res,
			ChatId: req.ChatId,
			Markup: keyboardMiddleware(user),
		},
		&userId,
	)
}

func (c *contextHandler) begin(ctx context.Context, userId int, req domain.Request) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	user, e := c.service.GetUser(ctx, userId)
	if e != nil {
		return newErrResponse(e, userId)
	}

	res, e := c.service.BeginStory(ctx, userId, &user)
	if e != nil {
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.Request{
			Data:   res,
			ChatId: req.ChatId,
			Markup: keyboardMiddleware(user),
		},
		&userId,
	)
}

func (c *contextHandler) editCharacter(ctx context.Context, userId int, req domain.Request) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	user, e := c.service.GetUser(ctx, userId)
	if e != nil {
		return newErrResponse(e, userId)
	}

	if req.Command != nil {
		keyboard := make(map[int][]domain.Button)
		keyboard[1] = []domain.Button{
			{
				Text: "Отменить",
				Data: string(domain.Cancel),
			},
		}

		return newSuccessResponse(
			domain.Request{
				Data:   "Отправьте новое описание персонажа",
				ChatId: req.ChatId,
				Markup: domain.InlineMarkup{
					Keyboard: keyboard,
				},
				ReplyTo: &req.MessageId,
			},
			nil,
		)
	}

	res, e := c.service.EditCharacter(ctx, userId, &user, req.Data)
	if e != nil {
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.Request{
			Data:   res,
			ChatId: req.ChatId,
			Markup: keyboardMiddleware(user),
		},
		&userId,
	)
}

func (c *contextHandler) editWorld(ctx context.Context, userId int, req domain.Request) response {
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, c.timeouts.RequestTimeout)
	defer cancel()

	user, e := c.service.GetUser(ctx, userId)
	if e != nil {
		return newErrResponse(e, userId)
	}

	if req.Command != nil {
		keyboard := make(map[int][]domain.Button)
		keyboard[1] = []domain.Button{
			{
				Text: "Отменить",
				Data: string(domain.Cancel),
			},
		}

		return newSuccessResponse(
			domain.Request{
				Data:   "Отправьте новое описание игрового мира",
				ChatId: req.ChatId,
				Markup: domain.InlineMarkup{
					Keyboard: keyboard,
				},
				ReplyTo: &req.MessageId,
			},
			nil,
		)
	}

	res, e := c.service.EditWorld(ctx, userId, &user, req.Data)
	if e != nil {
		return newErrResponse(e, userId)
	}

	return newSuccessResponse(
		domain.Request{
			Data:   res,
			ChatId: req.ChatId,
			Markup: keyboardMiddleware(user),
		},
		&userId,
	)
}
