package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func (s *server) handleIndex() http.HandlerFunc {
	type indexModel struct {
		Databases []string
	}

	return func (w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				println(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
			}

			name := r.Form.Get("dbname")
			pass := r.Form.Get("password")

			load := r.Form.Get("load") == "true"

			var repo ObjectRepository
			if load {
				data, err := Read(fmt.Sprintf("%s.db", name), []byte(pass))
				if err != nil {
					http.Error(w, "invalid database", http.StatusInternalServerError)
					return
				}

				repo, err = Load(data)
				if err != nil {
					http.Error(w, "invalid database", http.StatusInternalServerError)
					return
				}
			} else {
				repo = NewRepository()
			}

			s.dbName = name
			s.pass = []byte(pass)
			s.privateMode = false
			s.objects = repo

			http.Redirect(w, r, "/view", http.StatusSeeOther)
			return
		}

		t, err := template.ParseFiles("./static/layout.gohtml", "./static/index.gohtml")
		if err != nil {
			http.Error(w, "unable to parse templates", http.StatusInternalServerError)
			return
		}

		files, err := filepath.Glob("*.db")
		if err != nil {
			http.Error(w, "unable to detect databases", http.StatusInternalServerError)
			return
		}

		var names []string
		for _, f := range files {
			names = append(names, strings.TrimSuffix(f, ".db"))
		}

		err = t.ExecuteTemplate(w, "layout", &indexModel{
			Databases: names,
		})
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

func (s *server) handleView() http.HandlerFunc {
	type viewModel struct {
		TotalCount  int
		DbName      string
		Objects     []*Object
		PrivateMode bool
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				println(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
			}

			err = s.saveObject(r)
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
		for _, o := range s.objects.GetAll() {
			totalCount += o.Quantity
		}

		err = t.ExecuteTemplate(w, "layout", viewModel{
			TotalCount:  totalCount,
			DbName:      s.dbName,
			Objects:     s.objects.GetAll(),
			PrivateMode: s.privateMode,
		})
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

func (s *server) handleSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := Save(s.objects)
		if err != nil {
			println(err.Error())
			http.Error(w, "error", http.StatusInternalServerError)
		}

		err = Write(fmt.Sprintf("%s.db", s.dbName), s.pass, data)
		if err != nil {
			println(err.Error())
			http.Error(w, "error", http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/view", http.StatusSeeOther)
	}
}

func (s *server) handleClose() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.dbName = ""
		s.pass = nil
		s.privateMode = false
		s.objects = nil

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *server) handleEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			println(err.Error())
			http.Error(w, "bad request", http.StatusBadRequest)
		}

		id := r.Form.Get("id")
		object := s.objects.Get(id)
		if object == nil {
			http.Error(w, "object does not exist", http.StatusNotFound)
			return
		}

		if ! s.privateMode && object.Hidden {
			http.Error(w, "object can not be edited", http.StatusForbidden)
			return
		}

		if r.Method == http.MethodPost {
			err := s.saveObject(r)
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
}

func (s *server) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r * http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "unable to parse form", http.StatusBadRequest)
			return
		}
		err = s.objects.RemoveObject(r.Form.Get("id"))
		if err != nil {
			http.Error(w, "unable to remove object", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/view", http.StatusSeeOther)
	}
}

func (s *server) saveObject(r *http.Request) error {
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

		err = s.objects.PutObject(object)
		if err != nil {
			return err
		}
	}

	return nil
}
