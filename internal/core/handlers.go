package core

import (
	"fmt"
	"github.com/roelofruis/spullen"
	"github.com/roelofruis/spullen/internal/core/database"
	"github.com/roelofruis/spullen/internal/core/deletion"
	"github.com/roelofruis/spullen/internal/core/object"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) handleLoadDatabase(view *template.Template, isExistingDatabase bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := database.NewDatabaseForm(s.Finder)
		form.IsExistingDatabase = isExistingDatabase

		var alert = ""
		if r.Method == http.MethodPost {
			form.FillFromRequest(r)

			if form.Validate() {
				if s.Db.IsOpened() {
					err := s.Db.Close()
					if err != nil {
						log.Print(fmt.Sprintf("unable to close storage: %s", err.Error()))
						http.Error(w, "error", http.StatusInternalServerError)
						return
					}
				}

				err := s.Db.Open(form.Database, []byte(form.Password), form.IsExistingDatabase)
				if err == nil {
					s.PrivateMode = form.ParsedShowHiddenItems

					http.Redirect(w, r, "/view", http.StatusSeeOther)
					return
				}

				log.Printf("Error when trying to open storage: %s", err.Error())
				alert = "De database kon niet worden geopend. Het wachtwoord is fout of de database is corrupt."
			}
		}

		s.Render(w, view, &database.Database{
			Alert: alert,
			Form:  form,
		})
	}
}

func (s *Server) handleView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := object.EmptyForm()

		if r.Method == http.MethodPost {
			form.Id = s.MakeId()
			form.TimeAdded = strconv.FormatInt(time.Now().Truncate(time.Second).Unix(), 10)
			form.FillFromRequest(r)

			if form.Validate() {
				obj, err := form.GetObject()
				if err != nil {
					log.Print(err.Error())
					http.Error(w, "error", http.StatusInternalServerError)
					return
				}

				s.Objects.Put(obj)
				form = object.EmptyForm()
			}
		}

		s.Render(w, s.Views.View, &object.View{
			EditObject: object.EditObject{
				ExistingTags:         s.Objects.GetDistinctTags(s.PrivateMode),
				ExistingCategories:   s.Objects.GetDistinctCategories(s.PrivateMode),
				ExistingPropertyKeys: s.Objects.GetDistinctPropertyKeys(s.PrivateMode),
				Form:                 form,
			},
			DatabaseIsDirty: s.Db.IsDirty(),
			TotalCount:      s.Objects.Count(),
			DbName:          s.Db.Name(),
			Objects:         s.Objects.GetAll(),
			PrivateMode:     s.PrivateMode,
		})
	}
}

func (s *Server) handleEdit(o spullen.Object) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := object.FormFromObject(&o)

		var alert = ""
		if r.Method == http.MethodPost {
			form.FillFromRequest(r)

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

		s.Render(w, s.Views.Edit, &object.Edit{
			Alert: alert,
			EditObject: object.EditObject{
				ExistingTags:         s.Objects.GetDistinctTags(s.PrivateMode),
				ExistingCategories:   s.Objects.GetDistinctCategories(s.PrivateMode),
				ExistingPropertyKeys: s.Objects.GetDistinctPropertyKeys(s.PrivateMode),
				Form:                 form,
			},
		})
	}
}

func (s *Server) handleSplit(o spullen.Object) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if o.Quantity < 2 {
			http.Error(w, "object cannot be split", http.StatusConflict)
			return
		}

		form := object.FormFromObject(&o)
		form.Quantity = "1"

		var alert = ""
		if r.Method == http.MethodPost {
			form.Id = s.MakeId()
			form.TimeAdded = strconv.FormatInt(time.Now().Truncate(time.Second).Unix(), 10)
			form.FillFromRequest(r)

			if form.Validate() {
				splitObject, err := form.GetObject()
				if err != nil {
					alert = fmt.Sprintf("Error when getting object \n%s", err.Error())
				} else {
					o.Quantity -= splitObject.Quantity

					s.Objects.Put(splitObject)
					s.Objects.Put(&o)

					http.Redirect(w, r, "/view", http.StatusSeeOther)
					return
				}
			}
		}

		qty, err := strconv.ParseInt(form.Quantity, 10, 64)
		if err != nil {
			o.Quantity -= 1
		} else {
			o.Quantity -= int(qty)
		}

		original := object.FormFromObject(&o)

		s.Render(w, s.Views.Split, &object.Split{
			Alert: alert,
			EditObject: object.EditObject{
				ExistingTags:         s.Objects.GetDistinctTags(s.PrivateMode),
				ExistingCategories:   s.Objects.GetDistinctCategories(s.PrivateMode),
				ExistingPropertyKeys: s.Objects.GetDistinctPropertyKeys(s.PrivateMode),
				Form:                 form,
			},
			Original: original,
		})
	}
}

func (s *Server) handleDelete(o spullen.Object) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		original := object.FormFromObject(&o)

		var alert = ""
		form := &deletion.Form{Id: o.Id}
		if r.Method == http.MethodPost {
			form.RemovedAt = strconv.FormatInt(time.Now().Truncate(time.Second).Unix(), 10)
			form.Reason = r.PostFormValue("reason")

			if form.Validate() {
				// TODO: implement actual logic
				alert = "TODO: this is not implemented yet, object should now be deleted!"
			}
		}

		s.Render(w, s.Views.Delete, &deletion.Delete{
			Alert:    alert,
			Original: original,
			Form:     form,
		})
	}
}

func (s *Server) handleDestroy(object spullen.Object) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Objects.Remove(object.Id)

		http.Redirect(w, r, "/view", http.StatusSeeOther)
		return
	}
}

func (s *Server) handleSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.Db.Persist()
		if err != nil {
			log.Print(err.Error())
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
			log.Print(fmt.Sprintf("unable to persist storage: %s", err.Error()))
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		err = s.Db.Close()
		if err != nil {
			log.Print(fmt.Sprintf("unable to close storage: %s", err.Error()))
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
