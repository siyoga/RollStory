package router

import (
	"context"
	"fmt"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (r *router) requestWithMessageResponse(ctx context.Context, u domain.Update, route *Route) {
	r.lockRoute(u.Message.From.Id, route)
	res := route.handler(ctx, u.Message.From.Id, u.Message)
	r.logRequest(u, res.error)

	if res.release != nil {
		r.unlockRoute(*res.release)
	}

	var msg string
	if res.error != nil {
		if r.debug {
			msg = fmt.Sprintf("Error: %d, Reason: %s, Details: %s", res.code, res.error.Reason, res.error.Details.Error())
		} else {
			// TODO: maybe send error in private message if bot in prod mode
			msg = "Произошла ошибка. Попробуйте повторить запрос."
		}
	} else {
		msg = res.result.Message
	}

	if err := r.adpt.SendMessage(res.result.ChatId, msg); err != nil {
		r.logger.Error(errors.ErrTelegramSendMessage, err)
	}
}
