package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/roelofruis/spullen/internal/data"
)

func (app *application) handleListObjects(w http.ResponseWriter, r *http.Request) {
	query := contextGetQuery(r)

	objects, err := app.models.Objects.GetAll(query.Values.Get("name"))
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

	// TODO: validate

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
