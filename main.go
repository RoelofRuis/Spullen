package main

import (
	"fmt"
	"github.com/roelofruis/spullen/internal/database"
	"github.com/roelofruis/spullen/internal/spullen"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	devMode := os.Getenv("MODE") == "DEV"
	dbRoot := os.Getenv("DBROOT")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	repo := spullen.NewStorableObjectRepository()

	db := database.NewDatabase()
	db.Register(repo)

	server := &spullen.Server{
		Router: http.ServeMux{},
		Views:  &spullen.Views{},

		DevMode:     devMode,
		PrivateMode: true,

		Finder:  &spullen.Finder{Root: dbRoot},
		Db:      db,
		Objects: repo,
	}

	server.Templates()
	server.Routes()

	log.Printf("starting server on localhost:%s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), server)
	if err != nil {
		log.Fatal(err)
	}
}
