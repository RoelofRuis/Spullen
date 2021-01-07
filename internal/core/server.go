package core

import (
	"github.com/roelofruis/spullen"
	"html/template"
	"net/http"
)

type Server struct {
	Router http.ServeMux
	Views  *Views

	DevMode     bool
	PrivateMode bool

	Finder  *Finder
	Db      spullen.Database
	Objects spullen.ObjectRepository
}

type Views struct {
	Open  *template.Template
	New   *template.Template
	View  *template.Template
	Edit  *template.Template
	Delete *template.Template
	Split *template.Template
}

func (s *Server) Templates() {
	s.Views.Open = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/open.gohtml"))
	s.Views.New = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/new.gohtml"))

	s.Views.View = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/object-form.gohtml", "./static/view.gohtml"))
	s.Views.Edit = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/object-form.gohtml", "./static/edit.gohtml"))
	s.Views.Delete = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/object-original.gohtml", "./static/delete.gohtml"))
	s.Views.Split = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/object-form.gohtml", "./static/object-original.gohtml", "./static/split.gohtml"))
}

func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.handleOpen())
	s.Router.HandleFunc("/new", s.handleNew())

	s.Router.HandleFunc("/view", s.withDatabase(s.handleView()))
	s.Router.HandleFunc("/edit", s.withDatabase(s.handleEdit()))
	s.Router.HandleFunc("/split", s.withDatabase(s.handleSplit()))
	s.Router.HandleFunc("/delete", s.withDatabase(s.handleDelete()))
	s.Router.HandleFunc("/save", s.withDatabase(s.handleSave()))
	s.Router.HandleFunc("/close", s.withDatabase(s.handleClose()))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.DevMode {
		s.Templates()
	}

	s.Router.ServeHTTP(w, r)
}

func (s *Server) withDatabase(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Db.IsOpened() {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		h(w, r)
	}
}
