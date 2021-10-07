package main

import (
	"net/http"
)

func (app *application) handleListObjects(w http.ResponseWriter, r *http.Request) {
	objects, err := app.models.Objects.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"objects": objects}, nil)
}
