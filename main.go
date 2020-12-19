package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

type App struct {
	authenticated bool
	dbName        string
	pass          []byte
	privateMode   bool

	objects ObjectRepository
}

var app = &App {
	authenticated: false,
	dbName:        "",
	pass:          nil,
	privateMode:   false,

	objects: nil,
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/close", closeHandler)

	log.Print("started server on localhost:8080")

	http.ListenAndServe(":8080", nil)
}
