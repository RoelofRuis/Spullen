package spullen

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var index = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/index.gohtml"))
var view = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/view.gohtml"))
var edit = template.Must(template.ParseFiles("./static/layout.gohtml", "./static/edit.gohtml"))

func (s *Server) handleIndex() http.HandlerFunc {
	type indexModel struct {
		Alert string

		Databases []string
		Form      *IndexForm
	}

	return func(w http.ResponseWriter, r *http.Request) {
		form := &IndexForm{}

		var loadingAlert = ""
		if r.Method == http.MethodPost {
			form.ExistingDatabaseName = r.PostFormValue("existing-db")
			form.NewDatabaseName = r.PostFormValue("new-db")
			form.Password = r.PostFormValue("password")
			form.PrivateMode = r.PostFormValue("private-mode")

			if form.Validate() {
				if s.Db.IsOpened() {
					s.Db.Close()
				}

				repo, err := s.Db.Open(form.database, []byte(form.Password), !form.isNew)
				if err == nil {
					s.PrivateMode = form.isPrivateMode
					s.Objects = repo

					http.Redirect(w, r, "/view", http.StatusSeeOther)
					return
				}

				loadingAlert = "De database kon niet worden geopend. Het wachtwoord is fout of de database is corrupt."
			}
		}

		names, err := s.Finder.FindDatabases()
		if err != nil {
			http.Error(w, "unable to detect databases", http.StatusInternalServerError)
			return
		}

		err = index.ExecuteTemplate(w, "layout", &indexModel{
			Alert:     loadingAlert,
			Databases: names,
			Form:      form,
		})
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

func (s *Server) handleView() http.HandlerFunc {
	type viewModel struct {
		Alert string

		TotalCount  int
		DbName      string
		Objects     []*Object
		PrivateMode bool

		Form *ObjectForm
	}

	return func(w http.ResponseWriter, r *http.Request) {
		form := EmptyForm()

		var alert = ""
		if r.Method == http.MethodPost {
			form.Id = randSeq(16)
			form.TimeAdded = strconv.FormatInt(time.Now().Truncate(time.Second).Unix(), 10)
			form.Name = r.PostFormValue("name")
			form.Quantity = r.PostFormValue("quantity")
			form.Categories = r.PostFormValue("categories")
			form.Tags = r.PostFormValue("tags")
			form.Properties = r.PostFormValue("properties")
			form.Hidden = r.PostFormValue("hidden")
			form.Notes = r.PostFormValue("notes")

			if form.Validate() {
				obj, err := form.GetObject()
				if err != nil {
					alert = fmt.Sprintf("Error when getting object from form\n%s", err.Error())
				}

				_ = s.Objects.PutObject(obj)
				form = EmptyForm()
			}
		}

		totalCount := 0
		for _, o := range s.Objects.GetAll() {
			totalCount += o.Quantity
		}

		err := view.ExecuteTemplate(w, "layout", viewModel{
			Alert:       alert,
			TotalCount:  totalCount,
			DbName:      s.Db.Name(),
			Objects:     s.Objects.GetAll(),
			PrivateMode: s.PrivateMode,
			Form:        form,
		})
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

func (s *Server) handleSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.Db.Persist()
		if err != nil {
			println(err.Error())
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/view", http.StatusSeeOther)
	}
}

func (s *Server) handleClose() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.Db.Persist()
		if err != nil {
			println(err.Error())
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		s.Db.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *Server) handleEdit() http.HandlerFunc {
	type EditModel struct {
		Alert string

		Form *ObjectForm
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			println(err.Error())
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		id := r.Form.Get("id")
		object := s.Objects.Get(id)
		if object == nil {
			http.Error(w, "object does not exist", http.StatusNotFound)
			return
		}

		if !s.PrivateMode && object.Hidden {
			http.Error(w, "object can not be edited", http.StatusForbidden)
			return
		}

		form := FormFromObject(object)

		var alert = ""
		if r.Method == http.MethodPost {
			form.Name = r.PostFormValue("name")
			form.Quantity = r.PostFormValue("quantity")
			form.Categories = r.PostFormValue("categories")
			form.Tags = r.PostFormValue("tags")
			form.Properties = r.PostFormValue("properties")
			form.Hidden = r.PostFormValue("hidden")
			form.Notes = r.PostFormValue("notes")

			if form.Validate() {
				obj, err := form.GetObject()
				if err != nil {
					alert = fmt.Sprintf("Error when getting object\n%s", err.Error())
				} else {
					_ = s.Objects.PutObject(obj)

					http.Redirect(w, r, "/view", http.StatusSeeOther)
					return
				}
			}
		}

		err = edit.ExecuteTemplate(w, "layout", EditModel{Form: form, Alert: alert})
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

func (s *Server) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "unable to parse form", http.StatusBadRequest)
			return
		}
		err = s.Objects.RemoveObject(r.Form.Get("id"))
		if err != nil {
			http.Error(w, "unable to remove object", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/view", http.StatusSeeOther)
	}
}
