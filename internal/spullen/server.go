package spullen

import (
	"html/template"
	"net/http"
)

type Server struct {
	Router http.ServeMux
	Views  *Views

	DevMode     bool
	PrivateMode bool

	Finder  *Finder
	Db      Database
	Objects ObjectRepository
}

type Views struct {
	Index *template.Template
	View  *template.Template
	Edit  *template.Template
	Split *template.Template
}

func (s *Server) Templates() {
	s.Views.Index = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/index.gohtml"))
	s.Views.View = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/view.gohtml"))
	s.Views.Edit = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/edit.gohtml"))
	s.Views.Split = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/split.gohtml"))
}

func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.handleIndex())
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
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}
