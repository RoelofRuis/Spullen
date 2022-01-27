package main

import (
	"github.com/roelofruis/spullen/internal/db"
	"github.com/roelofruis/spullen/internal/jsonlog"
	"github.com/roelofruis/spullen/internal/model"
	"os"
)

type application struct {
	logger *jsonlog.Logger
	models model.Model
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := application{
		logger: logger,
		models: model.NewModel(db.NewProxy()),
	}

	err := app.serve()
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}
}
