package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type IndexModel struct {
	TotalCount  int
	Objects     []*Object
	PrivateMode bool
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			println(err.Error())
			http.Error(w, "bad request", http.StatusBadRequest)
		}

		name := r.Form.Get("name")
		pass := r.Form.Get("password")

		repo := NewRepository()

		app = &App{
			authenticated: true,
			path:          fmt.Sprintf("%s.db", name),
			pass:          []byte(pass),
			privateMode:   false,
			objects:       repo,
		}

		http.Redirect(w, r, "/view", http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("./static/layout.gohtml", "./static/index.gohtml")
	if err != nil {
		http.Error(w, "unable to parse templates", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		fmt.Print(err.Error())
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if ! app.authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			println(err.Error())
			http.Error(w, "bad request", http.StatusBadRequest)
		}

		err = saveObject(r)
		if err != nil {
			println(err.Error())
			http.Error(w, "error", http.StatusInternalServerError)
		}
	}

	t, err := template.ParseFiles("./static/layout.gohtml", "./static/view.gohtml")
	if err != nil {
		http.Error(w, "unable to parse templates", http.StatusInternalServerError)
		return
	}

	totalCount := 0
	for _, o := range app.objects.GetAll() {
		totalCount += o.Quantity
	}

	err = t.ExecuteTemplate(w, "layout", IndexModel{
		TotalCount:  totalCount,
		Objects:     app.objects.GetAll(),
		PrivateMode: app.privateMode,
	})
	if err != nil {
		fmt.Print(err.Error())
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if ! app.authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data, err := Save(app.objects)
	if err != nil {
		println(err.Error())
		http.Error(w, "error", http.StatusInternalServerError)
	}

	err = Write(app.path, app.pass, data)
	if err != nil {
		println(err.Error())
		http.Error(w, "error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/view", http.StatusSeeOther)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	if ! app.authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		println(err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	id := r.Form.Get("id")
	object := app.objects.Get(id)
	if object == nil {
		http.Error(w, "object does not exist", http.StatusNotFound)
		return
	}

	if ! app.privateMode && object.Hidden {
		http.Error(w, "object can not be edited", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		err := saveObject(r)
		if err != nil {
			println(err.Error())
			http.Error(w, "error", http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/view", http.StatusSeeOther)
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
	if ! app.authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}
	err = app.objects.RemoveObject(r.Form.Get("id"))
	if err != nil {
		http.Error(w, "unable to remove object", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view", http.StatusSeeOther)
}

func saveObject(r *http.Request) error {
	if len(r.PostForm.Get("name")) > 0 {
		object, err := ParseObjectForm(&ObjectForm{
			Id:         r.Form.Get("id"),
			Name:       r.PostForm.Get("name"),
			Quantity:   r.PostForm.Get("quantity"),
			Categories: r.PostForm.Get("categories"),
			Tags:       r.PostForm.Get("tags"),
			Properties: r.PostForm.Get("properties"),
			Hidden:     r.PostForm.Get("hidden"),
			Notes:      r.PostForm.Get("notes"),
		})
		if err != nil {
			return err
		}

		err = app.objects.PutObject(object)
		if err != nil {
			return err
		}
	}

	return nil
}
