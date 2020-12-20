package main

import "net/http"

type server struct {
	router http.ServeMux

	authenticated bool
	dbName string
	pass []byte
	privateMode bool

	objects ObjectRepository
}

func (s *server) routes() {
	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/edit", s.onlyAuthenticated(s.handleEdit()))
	s.router.HandleFunc("/view", s.onlyAuthenticated(s.handleView()))
	s.router.HandleFunc("/delete", s.onlyAuthenticated(s.handleDelete()))
	s.router.HandleFunc("/save", s.onlyAuthenticated(s.handleSave()))
	s.router.HandleFunc("/close", s.onlyAuthenticated(s.handleClose()))
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) onlyAuthenticated(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ! s.authenticated {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}