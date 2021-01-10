package core

import (
	"fmt"
	"github.com/roelofruis/spullen"
	"github.com/roelofruis/spullen/internal/util"
	"io"
	"log"
	"net/http"
)

func NewServer() *Server {
	return &Server{
		router:    http.ServeMux{},
		templates: &Templates{},
	}
}

type Server struct {
	router    http.ServeMux
	templates *Templates

	DevMode     bool
	PrivateMode bool

	Finder  *util.Finder
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

func (s *Server) Templates() {
	s.templates.Register("open", "./static/layout.gohtml", "./static/open.gohtml")
	s.templates.Register("new", "./static/layout.gohtml", "./static/new.gohtml")

	s.templates.Register("view", "./static/layout.gohtml", "./static/object-edit.gohtml", "./static/view.gohtml")
	s.templates.Register("edit", "./static/layout.gohtml", "./static/object-edit.gohtml", "./static/edit.gohtml")
	s.templates.Register("delete", "./static/layout.gohtml", "./static/object-display.gohtml", "./static/delete.gohtml")
	s.templates.Register("split", "./static/layout.gohtml", "./static/object-edit.gohtml", "./static/object-display.gohtml", "./static/split.gohtml")
}

func (s *Server) Routes() {
	s.router.HandleFunc("/", s.handleLoadDatabase("open", true))
	s.router.HandleFunc("/new", s.handleLoadDatabase("new", false))
	s.router.HandleFunc("/save", s.withDatabase(s.handleSave()))
	s.router.HandleFunc("/close", s.withDatabase(s.handleClose()))

	s.router.HandleFunc("/view", s.withDatabase(s.handleView()))
	s.router.HandleFunc("/edit", s.withDatabase(s.withValidObject(s.handleEdit)))
	s.router.HandleFunc("/split", s.withDatabase(s.withValidObject(s.handleSplit)))
	s.router.HandleFunc("/delete", s.withDatabase(s.withValidObject(s.handleDelete)))
	s.router.HandleFunc("/destroy", s.withDatabase(s.withValidObject(s.handleDestroy)))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.DevMode {
		s.Templates()
	}

	s.router.ServeHTTP(w, r)
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

func (s *Server) withValidObject(f func(o spullen.Object) http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Print(fmt.Sprintf("error when parsing form: %s", err.Error()))
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		id := spullen.ObjectId(r.Form.Get("id"))
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
		id = spullen.ObjectId(util.RandSeq(16))
		if !s.Objects.Has(id) {
			return id
		}
	}
}

type view struct {
	DevMode bool
	DbOpen  bool
	Version Version

	Data interface{}
}

func (s *Server) Render(w io.Writer, template string, data interface{}) {
	t, err := s.templates.Get(template)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = t.ExecuteTemplate(w, "layout", &view{
		DevMode: s.DevMode,
		DbOpen:  s.Db.IsOpened(),
		Version: s.Version,
		Data:    data,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
