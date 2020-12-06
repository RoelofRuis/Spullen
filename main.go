package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

func main() {
	o, err := LoadObjectList()
	if err != nil {
		log.Fatal(err)
	}
	obj := &Object{Name: "Stoel", Added: time.Now(), Tags: []string{"bats", "knats"}}
	o.AddObject(obj)
	o.Save()
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/index.html")

	o, _ := LoadObjectList()

	t.Execute(w, o)
}