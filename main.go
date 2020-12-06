package main

import (
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/index.html")

	o, _ := LoadObjectList()

	t.Execute(w, o)
}