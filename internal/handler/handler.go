package handler

import "github.com/siyoga/rollstory/internal/handler/router"

type (
	Handler interface {
		fillHandlers(r router.Router)
	}
)
