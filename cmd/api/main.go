package main

import (
	"log"
	"os"
)

type application struct {
	logger *log.Logger
}

func main() {
	app := application{
		logger: log.New(os.Stdout, "", log.Ltime),
	}

	app.logger.Printf("Starting API")

	err := app.serve()
	if err != nil {
		app.logger.Fatal(err)
	}
}
