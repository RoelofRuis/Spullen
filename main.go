package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var o Storage

var privateMode = true

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	storage, err := NewFileStorage()
	if err != nil {
		log.Fatal(err)
	}
	o = storage

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/delete", deleteHandler)

	fmt.Print("started server on localhost:8080")

	http.ListenAndServe(":8080", nil)
}

type IndexModel struct {
	Objects     *map[string]*Object
	PrivateMode bool
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "unable to parse form", http.StatusBadRequest)
			return
		}

		if len(r.PostForm.Get("name")) > 0 {
			object, err := ParseObjectForm(&ObjectForm{
				Name: r.PostForm.Get("name"),
				Quantity: r.PostForm.Get("quantity"),
				Categories: r.PostForm.Get("categories"),
				Tags: r.PostForm.Get("tags"),
				Properties: r.PostForm.Get("properties"),
				Private: r.PostForm.Get("private"),
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = o.AddObject(object)
			if err != nil {
				http.Error(w, "unable to add object", http.StatusInternalServerError)
				return
			}
		}
	}

	t, err := template.ParseFiles("./static/layout.gohtml", "./static/index.gohtml")
	if err != nil {
		http.Error(w, "unable to parse templates", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "layout", IndexModel{
		Objects:     o.GetAll(),
		PrivateMode: privateMode,
	})
	if err != nil {
		fmt.Print(err.Error())
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}
	id := r.Form.Get("id")
	object := o.Get(id)
	if object == nil {
		http.Error(w, "object does not exist", http.StatusNotFound)
		return
	}
	if privateMode && object.Private {
		http.Error(w, "object can not be edited", http.StatusForbidden)
		return
	}

	t, err := template.ParseFiles("./static/layout.gohtml", "./static/edit.gohtml")
	if err != nil {
		http.Error(w, "unable to parse templates", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "layout", MakeForm(object))
	if err != nil {
		fmt.Print(err.Error())
	}
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
