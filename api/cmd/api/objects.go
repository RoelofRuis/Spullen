package main

import (
	"errors"
	"github.com/roelofruis/spullen/internal/db"
	"github.com/roelofruis/spullen/internal/model"
	"github.com/roelofruis/spullen/internal/request"
	"github.com/roelofruis/spullen/internal/validator"
	"net/http"
)

func (app *application) handleListObjects(w http.ResponseWriter, r *http.Request) {
	var input struct {
		request.Filters
		Name string
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = request.ReadString(qs, "name", "")
	input.Filters.Page = request.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = request.ReadInt(qs, "page_size", 20, v)
	input.Filters.Sort = request.ReadString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id"}

	if request.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	objects, err := app.models.Objects.GetAll(input.Name, 0)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrNoDataSource):
			app.unauthorizedResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"objects": objects}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) handleCreateObject(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Quantity    int    `json:"quantity"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	object := model.NewObject(input.Name, input.Description)

	if input.Quantity > 0 {
		object.ChangeQuantity(input.Quantity, "")
	}

	v := validator.New()
	if object.Validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Objects.Insert(object); err != nil {
		switch {
		case errors.Is(err, db.ErrNoDataSource):
			app.unauthorizedResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"object": object}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) handleTagObject(w http.ResponseWriter, r *http.Request) {
	objectId, err := request.ReadIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		TagId    int `json:"tag_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	object, err := app.models.Objects.GetOne(model.ObjectID(objectId))
	if err != nil {
		switch {
		case errors.Is(err, db.ErrNoSuchRecord):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	tag, err := app.models.Tags.GetOne(model.TagID(input.TagId))
	if err != nil {
		switch {
		case errors.Is(err, db.ErrNoSuchRecord):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	object.AttachTag(tag)
	if err := app.models.Objects.Insert(object); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"object": object}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
