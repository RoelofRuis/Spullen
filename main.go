package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

type App struct {
	authenticated bool
	path string
	pass []byte
	privateMode bool

	objects ObjectRepository
}

var app = &App {
	authenticated: false,
	path: "",
	pass: nil,
	privateMode: false,

	objects: nil,
}

var privateMode = true

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/delete", deleteHandler)

	log.Print("started server on localhost:8080")

	http.ListenAndServe(":8080", nil)
}
