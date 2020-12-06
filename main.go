package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

var o *ObjectList

func main() {
	objectList, err := LoadObjectList()
	if err != nil {
		log.Fatal(err)
	}
	o = objectList

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/insert", insertHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/index.html")

	t.Execute(w, o)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}

	o.AddObject(&Object{
		Name: r.PostForm.Get("name"),
		Added: time.Now().Truncate(time.Second),
	})
	o.Save()

	http.Redirect(w, r, "/", http.StatusFound)
}