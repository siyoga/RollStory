package router

import (
	"context"
	"sync"

	"github.com/siyoga/rollstory/internal/adapter/telegram"
	"github.com/siyoga/rollstory/internal/errors"
	"github.com/siyoga/rollstory/internal/logger"
)

type (
	Router interface {
		Handle(cmd string, handler cmdHandler) *Route
		DefaultHandle(handler cmdHandler) *Route
		Run()
		Stop()
	}

	Handler interface {
		FillHandlers(r Router)
	}

	router struct {
		debug bool

		mu           sync.Mutex
		shutdownChan chan struct{}
		routes       map[string]*Route
		routesExec   map[int64]*Route // mapping to already executing routes by user
		defaultRoute *Route

		client telegram.Adapter

		ctxHandler  Handler
		gameHandler Handler

		logger logger.Logger
	}
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
		routes:       make(map[string]*Route),
		routesExec:   make(map[int64]*Route),

		client: adapter,

		ctxHandler:  ctxHandler,
		gameHandler: gameHandler,

		logger: logger,
	}
}

func (r *Route) AddButton(btns ...button) {
	// rewrite slice to make it exact needed len and cap
	keyboard := make([]button, len(btns))

	for i, btn := range btns {
		keyboard[i] = btn
	}

	r.buttons = keyboard
}

func (r *router) Handle(cmd string, handler cmdHandler) *Route {
	route := &Route{
		name:    cmd,
		handler: handler,
	}

	r.routes[route.name] = route

	return route
}

func (r *router) DefaultHandle(handler cmdHandler) *Route {
	route := &Route{
		name:    "text",
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
				if upd.Message == nil {
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
