package main

import (
	"github.com/roelofruis/spullen/internal_/jsonlog"
	"os"

	"github.com/roelofruis/spullen/internal_/data"
)

type application struct {
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := application{
		logger: logger,
		models: data.NewModels(data.NewDBProxy()),
	}

	err := app.serve()
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
}
