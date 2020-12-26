package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	dbRoot := os.Getenv("DBROOT")

	server := &server{
		router: http.ServeMux{},

		privateMode: true,

		finder:  &Finder{root: dbRoot},
		storage: nil,
		objects: nil,
	}

	server.routes()

	log.Print("started server on localhost:8080")

	err := http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatal(err)
	}
}
