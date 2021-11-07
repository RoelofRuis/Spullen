package main

import (
	"errors"
	"net/http"

	"github.com/roelofruis/spullen/internal_/data"
)

func (app *application) handleListObjects(w http.ResponseWriter, r *http.Request) {
	objects, err := app.models.Objects.GetAll()
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
