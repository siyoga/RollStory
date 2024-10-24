package router

import (
	"context"
	"fmt"

	"github.com/siyoga/rollstory/internal/adapter/telegram"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (r *router) handleUpdate(ctx context.Context, u telegram.Update) {
	cmdName := u.Message.GetCommand()

	executingRoute := r.getExecRoute(u.Message.From.Id)
	if executingRoute != nil {
		r.requestWithMessageResponse(ctx, u, executingRoute)
	} else {
		if u.Message.IsCommand() {
			if _, ok := r.routes[cmdName]; ok {
				r.requestWithMessageResponse(ctx, u, r.routes[cmdName])
			} else {
				msg := domain.Message{
					Text:   "Команда не найдена",
					ChatId: u.Message.Chat.ID,
				}.ToRequest(nil)

				if err := r.client.SendMessage(msg); err != nil {
					r.logger.Error(errors.ErrTelegramSendMessage, errors.ErrAuthNumberAssignmentFailed)
				}
			}
		} else {
			r.requestWithMessageResponse(ctx, u, r.defaultRoute)
		}
	}
}

func (r *router) requestWithMessageResponse(ctx context.Context, u telegram.Update, route *Route) {
	r.lockRoute(u.Message.From.Id, route)
	res := route.handler(ctx, u.Message.From.Id, domain.Message{}.FromUpdate(u))
	r.logRequest(u, res.error)

	if res.release != nil {
		r.unlockRoute(*res.release)
	}

	if res.error != nil {
		if r.debug {
			res.result.Text = fmt.Sprintf("Error: %d, Reason: %s, Details: %s", res.code, res.error.Reason, res.error.Details.Error())
		} else {
			// TODO: maybe send error in private message if bot in prod mode
			res.result.Text = "Произошла ошибка. Попробуйте повторить запрос."
		}
	}

	keyboard := make([]telegram.Button, len(route.buttons))
	for i, btn := range route.buttons {
		keyboard[i] = telegram.Button{
			Text: btn.text,
		}
	}

	if err := r.client.SendMessage(res.result.ToRequest(keyboard)); err != nil {
		r.logger.Error(errors.ErrTelegramSendMessage, err)
	}
}
