package core

import (
	"fmt"
	"github.com/roelofruis/spullen"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	Router http.ServeMux
	Views  *Views

	DevMode     bool
	PrivateMode bool

	Finder  *Finder
	Db      spullen.Database
	Objects spullen.ObjectRepository

	Version Version
}

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
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
	s.Router.HandleFunc("/edit/", s.withDatabase(s.withValidObject(s.handleEdit)))
	s.Router.HandleFunc("/split/", s.withDatabase(s.withValidObject(s.handleSplit)))
	s.Router.HandleFunc("/delete/", s.withDatabase(s.withValidObject(s.handleDelete)))
	s.Router.HandleFunc("/destroy/", s.withDatabase(s.withValidObject(s.handleDestroy)))
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

func (s *Server) withValidObject(f func(object spullen.Object) http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.Path, "/")

		if len(urlParts) < 2 {
			log.Print(fmt.Sprintf("unable to detect id"))
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}

		id := spullen.ObjectId(urlParts[2])

		err := r.ParseForm()
		if err != nil {
			log.Print(fmt.Sprintf("error when parsing form: %s", err.Error()))
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		object := s.Objects.Get(id)
		if object == nil {
			http.Error(w, "object does not exist", http.StatusNotFound)
			return
		}

		if !s.PrivateMode && object.Hidden {
			http.Error(w, "object can not be edited", http.StatusForbidden)
			return
		}

		f(*object)(w, r)
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
		Version: s.Version,
	}
}
