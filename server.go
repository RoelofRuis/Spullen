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
	s.router.HandleFunc("/edit", s.handleEdit())
	s.router.HandleFunc("/view", s.handleView())
	s.router.HandleFunc("/delete", s.handleDelete())
	s.router.HandleFunc("/save", s.handleSave())
	s.router.HandleFunc("/close", s.handleClose())
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}