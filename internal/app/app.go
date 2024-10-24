package app

import (
	"github.com/siyoga/rollstory/internal/app/deps"
)

type (
	Application interface {
		Run()
	}

	application struct {
		deps deps.Dependencies
	}
)

func NewApp(cfgPath string) Application {
	deps, err := deps.NewDependencies(cfgPath)
	if err != nil {
		panic(err)
	}

	return &application{
		deps: deps,
	}
}

func (app *application) Run() {
	router := app.deps.Router()
	router.Run()

	app.deps.WaitForInterrupt()

	app.deps.Close()
}
