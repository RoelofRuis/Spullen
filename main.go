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

	dbRoot := os.Getenv("DBROOT")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	factory := spullen.NewObjectRepositoryFactory()

	server := &spullen.Server{
		Router: http.ServeMux{},

		PrivateMode: true,
		Mode: spullen.ModeUseEncryption | spullen.ModeUseGzip,

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
