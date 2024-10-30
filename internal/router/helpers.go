package router

import (
	"fmt"
	"slices"

	uuid "github.com/satori/go.uuid"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (r *router) initRoutes() {
	r.ctxHandler.FillHandlers(r)
	r.gameHandler.FillHandlers(r)
}

// TODO: make this as part of request handling
func keyboardMiddleware(user domain.UserInfo) domain.ReplyMarkup {
	rm := domain.ReplyMarkup{
		Keyboard:     make(map[int][]domain.Button, 2),
		IsPersistent: false,
		Resize:       true,
		OneTime:      true,
	}

	fmt.Println(user)

	if !user.IsStarted {
		rm.AddRow([]domain.Button{
			{
				Text: string(CharacterPlain),
			},
			{
				Text: string(WorldPlain),
			},
		})

		if user.World != "" && user.Character != "" {
			rm.AddRow([]domain.Button{
				{
					Text: string(BeginPlain),
				},
			})
		}
	} else {
		rm.AddRow([]domain.Button{
			{
				Text: string(NewGamePlain),
			},
		})
	}

	fmt.Println(rm)

	return rm
}

func (r *execRoute) isLinked(cmdName Command) bool {
	return slices.Contains(r.route.linked, cmdName)
}

func (r *router) getRoute(id uuid.UUID) *route {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.routes[id]
}

func (r *router) getCommandId(cmdName Command) (uuid.UUID, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id, ok := r.triggers[cmdName]
	return id, ok
}

func (r *router) getExecRoute(userId int) *execRoute {
	r.mu.Lock()
	defer r.mu.Unlock()

	execRoute, ok := r.routesExec[userId]
	if !ok {
		return nil
	}

	return execRoute
}

func (r *router) lockRoute(userId int, newExec execRoute) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.routesExec[userId]; !ok {
		r.routesExec[userId] = &newExec
	}
}

func (r *router) unlockRoute(userId int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.routesExec, userId)
}

func newSuccessResponse(res domain.Request, release *int) response {
	return response{
		result:  res,
		error:   nil,
		release: release,
	}
}

func newErrResponse(e *errors.Error, release int) response {
	return response{
		result:  domain.Request{},
		error:   e,
		release: &release,
	}
}
