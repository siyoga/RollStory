package router

import (
	"context"
	"sync"

	uuid "github.com/satori/go.uuid"
	"github.com/siyoga/rollstory/internal/adapter/telegram"
	"github.com/siyoga/rollstory/internal/errors"
	"github.com/siyoga/rollstory/internal/logger"
)

func New(
	isDebug bool,
	logger logger.Logger,

	adapter telegram.Adapter,

	ctxHandler Handler,
	gameHandler Handler,
) Router {
	return &router{
		debug: isDebug,

		mu:           sync.Mutex{},
		shutdownChan: make(chan struct{}),

		triggers:   make(map[Command]uuid.UUID),
		routes:     make(map[uuid.UUID]*route),
		routesExec: make(map[int]*execRoute),

		client: adapter,

		ctxHandler:  ctxHandler,
		gameHandler: gameHandler,

		logger: logger,
	}
}

func (r *route) Handle(handler cmdHandler) *route {
	r.handler = handler

	return r
}

// attach route to route passed as argument
func (r *route) LinkTo(routes ...Command) {
	r.linked = append(r.linked, routes...)
}

func (r *router) Route(triggers ...Command) *route {
	id := uuid.NewV4()
	route := &route{
		id: id,
	}

	for _, trigger := range triggers {
		r.triggers[trigger] = id
	}

	r.routes[id] = route
	return route
}

func (r *router) DefaultRoute(handler cmdHandler) *route {
	route := &route{
		id:      uuid.NewV4(),
		handler: handler,
	}

	r.defaultRoute = route

	return route
}

func (r *router) Run() {
	r.initRoutes()

	if r.defaultRoute == nil {
		r.logger.Panic("Please, provide default route", errors.ErrRouterNoDefaultRoute)
	}

	go func() {
		ctx := context.Background()
		updChan := r.client.Updates()

	Loop:
		for {
			select {
			case upd := <-updChan:
				if upd.Message == nil && upd.Callback == nil {
					continue
				}

				go r.handleUpdate(ctx, upd)

			case <-r.shutdownChan:
				break Loop
			}
		}
	}()
}

func (r *router) Stop() {
	r.shutdownChan <- struct{}{}
}
