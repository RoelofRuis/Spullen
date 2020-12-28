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

	mode := os.Getenv("MODE")
	dbRoot := os.Getenv("DBROOT")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbMode := spullen.ModeUseEncryption | spullen.ModeUseGzip
	if mode == "DEV" {
		dbMode = 0x0
	}

	factory := spullen.NewObjectRepositoryFactory()

	server := &spullen.Server{
		Router: http.ServeMux{},

		PrivateMode: true,
		DbMode:      dbMode,

		Finder:  &spullen.Finder{Root: dbRoot},
		Db: database.NewDatabase(factory),
	}

	server.Routes()

	log.Printf("starting server on localhost:%s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), server)
	if err != nil {
		log.Fatal(err)
	}
}
