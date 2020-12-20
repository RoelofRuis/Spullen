package main

import "net/http"

type server struct {
	router http.ServeMux

	privateMode bool

	finder *Finder
	storage Storage
	objects ObjectRepository
}

func (s *server) routes() {
	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/edit", s.withLoadedObjects(s.handleEdit()))
	s.router.HandleFunc("/view", s.withLoadedObjects(s.handleView()))
	s.router.HandleFunc("/delete", s.withLoadedObjects(s.handleDelete()))
	s.router.HandleFunc("/save", s.withLoadedObjects(s.handleSave()))
	s.router.HandleFunc("/close", s.withLoadedObjects(s.handleClose()))
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) withLoadedObjects(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.objects == nil {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}