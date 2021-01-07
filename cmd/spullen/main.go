package main

import (
	"fmt"
	"github.com/roelofruis/spullen/internal/core"
	"github.com/roelofruis/spullen/internal/database"
	"github.com/roelofruis/spullen/internal/repository"
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

	versionManager := repository.NewVersionManager(1)
	objectRepo := repository.NewStorableObjectRepository(&core.ObjectMarshallerImpl{})

	var db *database.FileDatabase
	if devMode {
		db = database.NewDatabase(false, false)
	} else {
		db = database.NewDatabase(true, true)
	}

	_ = db.Register("object-repository", objectRepo)
	_ = db.Register("version-manager", versionManager)

	server := &core.Server{
		Router: http.ServeMux{},
		Views:  &core.Views{},

		DevMode:     devMode,
		PrivateMode: true,

		Finder:  &core.Finder{Root: dbRoot},
		Db:      db,
		Objects: objectRepo,
		Version: versionManager,
	}

	server.Templates()
	server.Routes()

	log.Printf("starting server on localhost:%s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), server)
	if err != nil {
		log.Fatal(err)
	}
}
