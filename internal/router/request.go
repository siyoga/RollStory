package router

import (
	"context"
	"fmt"

	"github.com/siyoga/rollstory/internal/adapter/telegram"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (r *router) handleUpdate(ctx context.Context, u telegram.Update) {
	req := u.ToDomain()

	executingRoute := r.getExecRoute(req.From.Id)
	switch executingRoute {
	case nil:
		if req.Command != nil {
			if cmdId, ok := r.getCommandId(Command(*req.Command)); ok {
				r.requestWithMessageResponse(ctx, req, r.getRoute(cmdId))
			} else {
				r.commandNotFound(req.ChatId)
			}
		} else {
			r.requestWithMessageResponse(ctx, req, r.defaultRoute)
		}

	// cancel execution of the current command and start execution of a new called command
	default:
		route := executingRoute.route

		if req.Command != nil {
			cmd := Command(*req.Command)

			if cmd == CancelSlash {
				r.cancelExec(executingRoute)
				return
			}

			if !executingRoute.isLinked(cmd) {
				r.cancelExec(executingRoute)
			}

			cmdId, ok := r.getCommandId(cmd)
			if !ok {
				r.commandNotFound(req.ChatId)
				return
			}
			route = r.getRoute(cmdId)
		}

		r.clearMarkups(executingRoute)
		r.requestWithMessageResponse(ctx, req, route)
	}
}

func (r *router) cancelExec(route *execRoute) {
	r.unlockRoute(route.userId)

	markups := &domain.InlineMarkup{
		Keyboard: make(map[int][]domain.Button),
	}

	if err := r.client.EditMessage(route.respChatId, route.respMessageId, "Отменено", markups); err != nil {
		r.logger.Error(errors.ErrTelegramSendMessage, err)
	}
}

func (r *router) clearMarkups(route *execRoute) {
	r.unlockRoute(route.userId)

	if err := r.client.EditMessage(route.respChatId, route.respMessageId, route.respText, &domain.InlineMarkup{}); err != nil {
		r.logger.Error(errors.ErrTelegramSendMessage, err)
	}
}

func (r *router) commandNotFound(chatId int) {
	res := telegram.Response{}.FromDomain(domain.Request{
		Data:   "Команда не найдена",
		ChatId: chatId,
	})

	if _, err := r.client.SendMessage(res); err != nil {
		r.logger.Error(errors.ErrTelegramSendMessage, errors.ErrAuthNumberAssignmentFailed)
	}
}

// for now callbacks only used for cancel command executing
func (r *router) requestWithMessageResponse(ctx context.Context, req domain.Request, route *route) {
	res := route.handler(ctx, req.From.Id, req)
	r.logRequest(req, res.error)

	if res.error != nil {
		res.release = &req.From.Id
		res.result.ChatId = req.ChatId

		if r.debug {
			res.result.Data = fmt.Sprintf("ERROR | Reason: %s, Details: %s\n", res.error.Reason, res.error.Details.Error())
		} else {
			// TODO: maybe send error in private message if bot in prod mode
			res.result.Data = "Произошла ошибка. Попробуйте повторить запрос."
		}
	}

	msg, err := r.client.SendMessage(telegram.Response{}.FromDomain(res.result))
	if err != nil {
		r.logger.Error(errors.ErrTelegramSendMessage, err)
	}

	r.lockRoute(req.From.Id, execRoute{
		route:  route,
		userId: req.From.Id,

		respText:      msg.Text,
		respChatId:    msg.Chat.ID,
		respMessageId: msg.Id,
	})

	if res.release != nil {
		r.unlockRoute(*res.release)
	}
}
