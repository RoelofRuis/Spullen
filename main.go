package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	server := &server{
		router:        http.ServeMux{},
		dbName:        "",
		pass:          nil,
		privateMode:   false,
		objects:       nil,
	}

	server.routes()

	log.Print("started server on localhost:8080")

	err := http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatal(err)
	}
}
