package deps

import (
	"github.com/siyoga/rollstory/internal/router"
)

func (d *dependencies) Router() router.Router {
	if d.router == nil {
		d.router = router.New(d.cfg.Bot.Debug, d.log, d.TelegramAdapter(), d.ContextHandler(), d.GameHandler())

		d.closeCallbacks = append(d.closeCallbacks, func() {
			msg := "stop bot router"
			d.router.Stop()
			d.log.Zap().Info(msg)
		})
	}

	return d.router
}

func (d *dependencies) ContextHandler() router.Handler {
	if d.contextHandler == nil {
		d.contextHandler = router.NewContextHandler(d.cfg.Timeouts, d.ContextService())
	}

	return d.contextHandler
}

func (d *dependencies) GameHandler() router.Handler {
	if d.gameHandler == nil {
		d.gameHandler = router.NewGameHandler(d.cfg.Timeouts, d.GameService(), d.ContextService())
	}

	return d.gameHandler
}
