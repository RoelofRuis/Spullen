package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var o Storage

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	storage, err := NewFileStorage()
	if err != nil {
		log.Fatal(err)
	}
	o = storage

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/delete", deleteHandler)

	fmt.Print("started server on localhost:8080")

	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "unable to parse form", http.StatusBadRequest)
			return
		}

		name := r.PostForm.Get("name")
		if len(name) > 0 {
			err := o.AddObject(&Object{
				Id:         randSeq(16),
				Name:       strings.ToLower(name),
				Quantity:   1,
				Added:      time.Now().Truncate(time.Second),
				Categories: nil,
				Tags:       nil,
				Properties: nil,
				Private:    false,
			})
			if err != nil {
				print(err.Error())
				http.Error(w, "unable to add object", http.StatusInternalServerError)
				return
			}
		}
	}

	t, _ := template.ParseFiles("./static/index.html")

	t.Execute(w, o.GetAll())
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}
	err = o.RemoveObject(r.Form.Get("id"))
	if err != nil {
		http.Error(w, "unable to remove object", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
