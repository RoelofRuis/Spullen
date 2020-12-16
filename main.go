package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var o Storage

var privateMode = true

func main() {
	data := "iedereen is liefde"
	cipher, err := Encrypt([]byte("je moedr"), []byte(data))
	if err != nil {
		log.Fatal(err)
	}
	plain, err := Decrypt([]byte("je moedr"), cipher)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(plain))

	rand.Seed(time.Now().UTC().UnixNano())

	storage, err := NewFileStorage()
	if err != nil {
		log.Fatal(err)
	}
	o = storage

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/delete", deleteHandler)

	log.Print("started server on localhost:8080")

	http.ListenAndServe(":8080", nil)
}
