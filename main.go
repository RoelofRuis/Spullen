package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

var o ObjectRepository

var privateMode = true

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	storage, err := NewFileStorage()
	if err != nil {
		log.Fatal(err)
	}
	o = storage

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/delete", deleteHandler)

	log.Print("started server on localhost:8080")

	http.ListenAndServe(":8080", nil)
}
