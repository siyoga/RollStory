package deps

import (
	"github.com/siyoga/rollstory/internal/service"
	ctxSvc "github.com/siyoga/rollstory/internal/service/context"
	"github.com/siyoga/rollstory/internal/service/game"
)

func (d *dependencies) ContextService() service.ContextService {
	if d.contextService == nil {
		d.contextService = ctxSvc.NewContextService(d.log, d.GptAdapter(), d.StoryRepository(), d.ThreadRepository())
	}

	return d.contextService
}

func (d *dependencies) GameService() service.GameService {
	if d.gameService == nil {
		d.gameService = game.NewGameService(d.log, d.GptAdapter(), d.ThreadRepository())
	}

	return d.gameService
}
