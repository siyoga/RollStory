package router

import (
	"context"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (r *router) getUpdates() domain.UpdatesChan {
	ch := make(chan domain.Update, r.batchSize)

	go func() {
		for {
			select {
			case <-r.shutdownChan:
				close(ch)
				return
			default:
			}

			updates, err := r.adpt.Updates(r.offset, r.batchSize)
			if err != nil {
				r.logger.Error(errors.ErrTelegramGetUpdates, err)
				continue
			}

			for _, update := range updates {
				if update.ID >= r.offset {
					r.offset = update.ID + 1
					ch <- update
				}
			}
		}
	}()

	return ch
}

func (r *router) run() {
	ctx := context.Background()
	updChan := r.getUpdates()

Loop:
	for {
		select {
		case upd := <-updChan:
			if upd.Message == nil {
				continue
			}

			go r.handleUpdate(ctx, upd)

		case <-r.shutdownChan:
			break Loop
		}
	}
}

func (r *router) handleUpdate(ctx context.Context, u domain.Update) {
	cmdName := u.Message.GetCommand()

	executingRoute := r.getExecRoute(u.Message.From.Id)
	if executingRoute != nil {
		r.requestWithMessageResponse(ctx, u, executingRoute)
	} else {
		if u.Message.IsCommand() {
			if _, ok := r.routes[cmdName]; ok {
				r.requestWithMessageResponse(ctx, u, r.routes[cmdName])
			} else {
				if err := r.adpt.SendMessage(u.Message.From.Id, "Команда не найдена."); err != nil {
					r.logger.Error(errors.ErrTelegramSendMessage, errors.ErrAuthNumberAssignmentFailed)
				}
			}
		} else {
			r.requestWithMessageResponse(ctx, u, r.defaultRoute)
		}
	}
}

func (r *router) getExecRoute(userId int64) *Route {
	r.mu.Lock()
	defer r.mu.Unlock()

	route, ok := r.routesExec[userId]
	if !ok {
		return nil
	}

	return route
}

func (r *router) lockRoute(userId int64, route *Route) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.routesExec[userId]; !ok {
		r.routesExec[userId] = route
	}
}

func (r *router) unlockRoute(userId int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.routesExec, userId)
}
