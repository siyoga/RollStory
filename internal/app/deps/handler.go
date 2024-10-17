package deps

import (
	"github.com/siyoga/rollstory/internal/handler"
	"github.com/siyoga/rollstory/internal/handler/router"
)

func (d *dependencies) Router() router.Router {
	if d.router == nil {
		d.router = router.New(d.cfg.Bot, d.TelegramAdapter(), d.log)

		d.closeCallbacks = append(d.closeCallbacks, func() {
			msg := "stop bot router"
			d.router.Stop()
			d.log.Zap().Info(msg)
		})
	}

	return d.router
}

func (d *dependencies) ContextHandler() handler.Handler {
	if d.contextHandler == nil {
		d.contextHandler = handler.NewContextHandler(d.cfg.Timeouts, d.ContextService(), d.Router())
	}

	return d.contextHandler
}

func (d *dependencies) GameHandler() handler.Handler {
	if d.gameHandler == nil {
		d.gameHandler = handler.NewGameHandler(d.cfg.Timeouts, d.GameService(), d.Router())
	}

	return d.gameHandler
}
