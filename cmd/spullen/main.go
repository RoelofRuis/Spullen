package main

import (
	"fmt"
	"github.com/roelofruis/spullen/internal/core"
	"github.com/roelofruis/spullen/internal/core/deletion"
	"github.com/roelofruis/spullen/internal/core/object"
	"github.com/roelofruis/spullen/internal/storage"
	"github.com/roelofruis/spullen/internal/util"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var VERSION = core.Version{Major: 0, Minor: 9, Patch: 1}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	devMode := os.Getenv("MODE") == "DEV"
	dbRoot := os.Getenv("DBROOT")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	objectRepo := object.NewRepository()
	deletionRepo := deletion.NewRepository()

	var db *storage.FileDatabase
	if devMode {
		db = storage.NewDatabase(false, false)
	} else {
		db = storage.NewDatabase(true, true)
	}

	_ = db.Register("object-repository", objectRepo)
	_ = db.Register("deletion-repository", deletionRepo)

	server := core.NewServer()
	server.DevMode = devMode
	server.Finder = &util.Finder{Root: dbRoot}
	server.Db = db
	server.Objects = objectRepo
	server.Deletions = deletionRepo
	server.Version = VERSION
	server.ObjectViewer = core.NewObjectViewer(objectRepo, deletionRepo)

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
