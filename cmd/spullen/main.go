package main

import (
	"fmt"
	"github.com/roelofruis/spullen/internal/core"
	"github.com/roelofruis/spullen/internal/core/object"
	"github.com/roelofruis/spullen/internal/database"
	"github.com/roelofruis/spullen/internal/util"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var VERSION = core.Version{Major: 0, Minor: 7, Patch: 6}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	devMode := os.Getenv("MODE") == "DEV"
	dbRoot := os.Getenv("DBROOT")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	objectRepo := object.NewStorableObjectRepository(&object.MarshallerImpl{})

	var db *database.FileDatabase
	if devMode {
		db = database.NewDatabase(false, false)
	} else {
		db = database.NewDatabase(true, true)
	}

	_ = db.Register("object-repository", objectRepo)

	server := &core.Server{
		Router: http.ServeMux{},
		Views:  &core.Views{},

		DevMode:     devMode,
		PrivateMode: true,

		Finder:  &util.Finder{Root: dbRoot},
		Db:      db,
		Objects: objectRepo,

		Version: VERSION,
	}

	server.Templates()
	server.Routes()

	log.Printf("spullen app [%s]", VERSION.String())
	if devMode {
		log.Printf("development mode enabled")
	}
	log.Printf("starting server on localhost:%s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), server)
	if err != nil {
		log.Fatal(err)
	}
}
