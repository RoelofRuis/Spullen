package main

import (
	"fmt"
	"github.com/roelofruis/spullen/internal_/validator"
	"net/http"
	"regexp"
)

var (
	FileRX = regexp.MustCompile("^[a-zA-Z0-9_]*$")
)

func (app *application) handleNewDatabase(w http.ResponseWriter, r *http.Request) {
	var input struct {
		File string `json:"file"`
		Key  string `json:"key"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.File != "", "file", "must not be empty")
	v.Check(validator.Matches(input.File, FileRX), "file", "can only contain alphanumeric characters and underscore")
	v.Check(input.Key != "", "key", "must not be empty")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.DB.Open(input.File, input.Key)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "Created new database\n")
}
