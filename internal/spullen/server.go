package spullen

import (
	"net/http"
)

type Server struct {
	Router http.ServeMux

	PrivateMode bool

	Mode DatabaseMode

	Finder  *Finder
	Db Database
	Objects ObjectRepository
}

func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.handleIndex())
	s.Router.HandleFunc("/edit", s.withDatabase(s.handleEdit()))
	s.Router.HandleFunc("/view", s.withDatabase(s.handleView()))
	s.Router.HandleFunc("/delete", s.withDatabase(s.handleDelete()))
	s.Router.HandleFunc("/save", s.withDatabase(s.handleSave()))
	s.Router.HandleFunc("/close", s.withDatabase(s.handleClose()))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
