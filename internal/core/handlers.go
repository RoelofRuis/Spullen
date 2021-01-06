package core

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) handleIndex() http.HandlerFunc {
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
					_ = s.Db.Close()
				}

				openExisting := !form.isNew

				err := s.Db.Open(form.database, []byte(form.Password), openExisting)
				if err == nil {
					s.PrivateMode = form.isPrivateMode

					http.Redirect(w, r, "/view", http.StatusSeeOther)
					return
				}

				println(err.Error())
				loadingAlert = "De database kon niet worden geopend. Het wachtwoord is fout of de database is corrupt."
			}
		}

		names, err := s.Finder.FindDatabases()
		if err != nil {
			http.Error(w, "unable to detect databases", http.StatusInternalServerError)
			return
		}

		err = s.Views.Index.ExecuteTemplate(w, "layout", &Index{
			AppInfo:   AppInfo{s.DevMode, loadingAlert},
			Databases: names,
			Form:      form,
		})
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

func (s *Server) handleView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := EmptyForm()

		var alert = ""
		if r.Method == http.MethodPost {
			form.Id = s.makeId()
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

				s.Objects.Put(obj)
				form = EmptyForm()
			}
		}

		err := s.Views.View.ExecuteTemplate(w, "layout", View{
			AppInfo: AppInfo{DevMode: s.DevMode, Alert: alert},
			EditObject: EditObject{
				ExistingTags: s.Objects.GetDistinctTags(s.PrivateMode),
				ExistingCategories: s.Objects.GetDistinctCategories(s.PrivateMode),
				ExistingPropertyKeys: s.Objects.GetDistinctPropertyKeys(s.PrivateMode),
				Form: form,
			},
			DatabaseIsDirty: s.Db.IsDirty(),
			TotalCount:      s.Objects.Count(),
			DbName:          s.Db.Name(),
			Objects:         s.Objects.GetAll(),
			PrivateMode:     s.PrivateMode,
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

		_ = s.Db.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *Server) handleSplit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			println(err.Error())
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		id := r.Form.Get("id")
		objectPointer := s.Objects.Get(id)
		if objectPointer == nil {
			http.Error(w, "object does not exist", http.StatusNotFound)
			return
		}

		object := *objectPointer

		if object.Quantity < 2 {
			http.Error(w, "object cannot be split", http.StatusBadRequest)
			return
		}

		if !s.PrivateMode && object.Hidden {
			http.Error(w, "object can not be edited", http.StatusForbidden)
			return
		}

		form := FormFromObject(&object)
		object.Quantity -= 1
		original := FormFromObject(&object)

		var alert = ""
		if r.Method == http.MethodPost {
			form.Id = s.makeId()
			form.TimeAdded = strconv.FormatInt(time.Now().Truncate(time.Second).Unix(), 10)
			form.Name = r.PostFormValue("name")
			form.Quantity = "1"
			form.Categories = r.PostFormValue("categories")
			form.Tags = r.PostFormValue("tags")
			form.Properties = r.PostFormValue("properties")
			form.Hidden = r.PostFormValue("hidden")
			form.Notes = r.PostFormValue("notes")

			if form.Validate() {
				splitObject, err := form.GetObject()
				if err != nil {
					alert = fmt.Sprintf("Error when getting object \n%s", err.Error())
				} else {
					s.Objects.Put(splitObject)
					s.Objects.Put(&object)

					http.Redirect(w, r, "/view", http.StatusSeeOther)
					return
				}
			}
		}

		err = s.Views.Split.ExecuteTemplate(w, "layout", Split{
			AppInfo: AppInfo{DevMode: s.DevMode, Alert: alert},
			EditObject: EditObject{
				ExistingTags: s.Objects.GetDistinctTags(s.PrivateMode),
				ExistingCategories: s.Objects.GetDistinctCategories(s.PrivateMode),
				ExistingPropertyKeys: s.Objects.GetDistinctPropertyKeys(s.PrivateMode),
				Form: form,
			},
			Original: original,
		})
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

func (s *Server) handleEdit() http.HandlerFunc {
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
					s.Objects.Put(obj)

					http.Redirect(w, r, "/view", http.StatusSeeOther)
					return
				}
			}
		}

		err = s.Views.Edit.ExecuteTemplate(w, "layout", Edit{
			AppInfo: AppInfo{DevMode: s.DevMode, Alert: alert},
			EditObject: EditObject{
				ExistingTags: s.Objects.GetDistinctTags(s.PrivateMode),
				ExistingCategories: s.Objects.GetDistinctCategories(s.PrivateMode),
				ExistingPropertyKeys: s.Objects.GetDistinctPropertyKeys(s.PrivateMode),
				Form: form,
			},
		})
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
		s.Objects.Remove(r.Form.Get("id"))

		http.Redirect(w, r, "/view", http.StatusSeeOther)
	}
}

func (s *Server) makeId() string {
	var id string
	for {
		id = randSeq(16)
		if !s.Objects.Has(id) {
			return id
		}
	}
}
