package main

import (
	"fmt"
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

	server := &server{
		router: http.ServeMux{},

		privateMode: true,

		finder:  &Finder{root: dbRoot},
		storage: nil,
		objects: nil,
	}

	server.routes()

	log.Printf("starting server on localhost:%s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), server)
	if err != nil {
		log.Fatal(err)
	}
}
