package main

import (
	"flag"

	"github.com/siyoga/rollstory/internal/app"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "config file dir")
	flag.Parse()

	application := app.NewApp(path)
	application.Run()
}
