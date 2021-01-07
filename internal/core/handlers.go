package core

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) handleNew() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		form, err := NewDatabaseForm(s.Finder)
		if err != nil {
			log.Fatal("Error when creating form database form", err.Error())
		}

		var alert = ""
		if r.Method == http.MethodPost {
			form.Database = r.PostFormValue("database")
			form.Password = r.PostFormValue("password")
			form.ShowHiddenItems = r.PostFormValue("show-hidden-items")

			if form.Validate(false) {
				if s.Db.IsOpened() {
					_ = s.Db.Close()
				}

				err := s.Db.Open(form.Database, []byte(form.Password), false)
				if err == nil {
					s.PrivateMode = form.ParsedShowHiddenItems

					http.Redirect(w, r, "/view", http.StatusSeeOther)
					return
				}

				log.Printf("Error when trying to open database: %s", err.Error())
				alert = "De database kon niet worden geopend. Het wachtwoord is fout of de database is corrupt."
			}
		}

		err = s.Views.New.ExecuteTemplate(w, "layout", &Database{
			AppInfo:   AppInfo{s.DevMode, alert},
			Form:      form,
		})
		if err != nil {
			fmt.Print(err.Error())
		}

	}
}

func (s *Server) handleOpen() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form, err := NewDatabaseForm(s.Finder)
		if err != nil {
			log.Fatal("Error when creating form database form", err.Error())
		}

		var alert = ""
		if r.Method == http.MethodPost {
			form.Database = r.PostFormValue("database")
			form.Password = r.PostFormValue("password")
			form.ShowHiddenItems = r.PostFormValue("show-hidden-items")

			if form.Validate(true) {
				if s.Db.IsOpened() {
					_ = s.Db.Close()
				}

				err := s.Db.Open(form.Database, []byte(form.Password), true)
				if err == nil {
					s.PrivateMode = form.ParsedShowHiddenItems

					http.Redirect(w, r, "/view", http.StatusSeeOther)
					return
				}

				log.Printf("Error when trying to open database: %s", err.Error())
				alert = "De database kon niet worden geopend. Het wachtwoord is fout of de database is corrupt."
			}
		}

		err = s.Views.Open.ExecuteTemplate(w, "layout", &Database{
			AppInfo:   AppInfo{s.DevMode, alert},
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
		form.Quantity = "1"

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
				splitObject, err := form.GetObject()
				if err != nil {
					alert = fmt.Sprintf("Error when getting object \n%s", err.Error())
				} else {
					object.Quantity -= splitObject.Quantity

					s.Objects.Put(splitObject)
					s.Objects.Put(&object)

					http.Redirect(w, r, "/view", http.StatusSeeOther)
					return
				}
			}
		}

		qty, err := strconv.ParseInt(form.Quantity, 10, 64)
		if err != nil {
			object.Quantity -= 1
		} else {
			object.Quantity -= int(qty)
		}

		original := FormFromObject(&object)

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

		original := FormFromObject(object)

		var alert = ""
		form := &DeleteForm{Id: id}
		if r.Method == http.MethodPost {
			// TODO:
		}

		err = s.Views.Delete.ExecuteTemplate(w, "layout", Delete{
			AppInfo: AppInfo{DevMode: s.DevMode, Alert: alert},
			Original: original,
			Form: form,
		})
		if err != nil {
			fmt.Print(err.Error())
		}
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
