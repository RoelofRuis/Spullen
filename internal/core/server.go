package core

import (
	"fmt"
	"github.com/roelofruis/spullen"
	"html/template"
	"log"
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
	Open   *template.Template
	New    *template.Template
	View   *template.Template
	Edit   *template.Template
	Delete *template.Template
	Split  *template.Template
}

func (s *Server) Templates() {
	s.Views.Open = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/open.gohtml"))
	s.Views.New = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/new.gohtml"))

	s.Views.View = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/object-edit.gohtml", "./static/view.gohtml"))
	s.Views.Edit = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/object-edit.gohtml", "./static/edit.gohtml"))
	s.Views.Delete = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/object-display.gohtml", "./static/delete.gohtml"))
	s.Views.Split = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/object-edit.gohtml", "./static/object-display.gohtml", "./static/split.gohtml"))
}

func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.handleLoadDatabase(s.Views.Open, true))
	s.Router.HandleFunc("/new", s.handleLoadDatabase(s.Views.New, false))
	s.Router.HandleFunc("/save", s.withDatabase(s.handleSave()))
	s.Router.HandleFunc("/close", s.withDatabase(s.handleClose()))

	s.Router.HandleFunc("/view", s.withDatabase(s.handleView()))
	s.Router.HandleFunc("/edit", s.withDatabase(s.withParsedForm(s.handleEdit())))
	s.Router.HandleFunc("/split", s.withDatabase(s.withParsedForm(s.handleSplit())))
	s.Router.HandleFunc("/delete", s.withDatabase(s.withParsedForm(s.handleDelete())))
	s.Router.HandleFunc("/destroy", s.withDatabase(s.withParsedForm(s.handleDestroy())))
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

func (s *Server) withParsedForm(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Print(fmt.Sprintf("error when parsing form: %s", err.Error()))
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		h(w, r)
	}
}

func (s *Server) MakeId() spullen.ObjectId {
	var id spullen.ObjectId
	for {
		id = spullen.ObjectId(randSeq(16))
		if !s.Objects.Has(id) {
			return id
		}
	}
}

func (s *Server) AppInfo() AppInfo {
	return AppInfo{
		DevMode: s.DevMode,
		DbOpen:  s.Db.IsOpened(),
	}
}
