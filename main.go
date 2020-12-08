package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var o *ObjectList

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	objectList, err := LoadObjectList()
	if err != nil {
		log.Fatal(err)
	}
	o = objectList

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/index.html")

	t.Execute(w, o)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	if len(name) > 0 {
		o.AddObject(&Object{
			Id: randSeq(16),
			Name: strings.ToLower(name),
			Added: time.Now().Truncate(time.Second),
			Categories: nil,
			Tags: nil,
			Properties: nil,
			Private: false,
		})
		o.Save()
	}


	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}
	o.RemoveObject(r.Form.Get("id"))
	o.Save()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}