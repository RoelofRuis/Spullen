package main

import (
	"errors"
	"github.com/roelofruis/spullen/internal/request"
	"github.com/roelofruis/spullen/internal/validator"
	"net/http"
	"time"

	"github.com/roelofruis/spullen/internal/data"
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

	if request.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	objects, err := app.models.Objects.GetAll(input.Name)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoDataSource):
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

func (app *application) handleAddObject(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	object := &data.Object{
		Added:    time.Now(),
		Name:     input.Name,
		Quantity: input.Quantity,
	}

	v := validator.New()
	if data.ValidateObject(v, object); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Objects.Insert(object); err != nil {
		switch {
		case errors.Is(err, data.ErrNoDataSource):
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
