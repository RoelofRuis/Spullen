package main

import (
	"errors"
	"github.com/roelofruis/spullen/internal/db"
	"github.com/roelofruis/spullen/internal/model"
	"github.com/roelofruis/spullen/internal/validator"
	"net/http"
)

func (app *application) handleListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := app.models.Tags.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, db.ErrNoDataSource):
			app.unauthorizedResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tags": tags}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) handleCreateTag(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tag := &model.Tag{
		Name:        input.Name,
		Description: input.Description,
		IsSystemTag: false,
	}

	v := validator.New()
	if model.ValidateTag(v, tag); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Tags.Insert(tag); err != nil {
		switch {
		case errors.Is(err, db.ErrNoDataSource):
			app.unauthorizedResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"tag": tag}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
