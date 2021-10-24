package main

import (
	"log"
	"os"

	"github.com/roelofruis/spullen/internal_/data"
)

type application struct {
	logger *log.Logger
	models data.Models
}

func main() {
	app := application{
		logger: log.New(os.Stdout, "", log.Ltime),
		models: data.NewModels(&data.DBProxy{}),
	}

	app.logger.Printf("Starting API")

	err := app.serve()
	if err != nil {
		app.logger.Fatal(err)
	}
}
