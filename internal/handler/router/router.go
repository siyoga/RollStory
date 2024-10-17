package router

import (
	"sync"

	"github.com/siyoga/rollstory/internal/adapter"
	"github.com/siyoga/rollstory/internal/config"
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

	router struct {
		debug     bool
		offset    int
		batchSize int

		mu           sync.Mutex
		shutdownChan chan struct{}
		routes       map[string]*Route
		routesExec   map[int64]*Route // mapping to already executing routes by user
		defaultRoute *Route

		adpt adapter.TelegramAdapter

		logger logger.Logger
	}
)

func New(
	cfg config.Bot,
	adapter adapter.TelegramAdapter,
	logger logger.Logger,
) Router {
	return &router{
		debug:     cfg.Debug,
		offset:    cfg.Offset,
		batchSize: cfg.BatchSize,

		mu:           sync.Mutex{},
		shutdownChan: make(chan struct{}),
		routes:       make(map[string]*Route),
		routesExec:   make(map[int64]*Route),

		adpt: adapter,

		logger: logger,
	}
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
	if r.defaultRoute == nil {
		r.logger.Panic("Please, provide default route", errors.ErrRouterNoDefaultRoute)
	}

	go r.run()
}

func (r *router) Stop() {
	r.shutdownChan <- struct{}{}
}
