package main

import (
	"errors"
	"fmt"
	"github.com/roelofruis/spullen/internal/data"
	"github.com/roelofruis/spullen/internal/validator"
	"net/http"
	"os"
	"path"
)

func (app *application) handleOpenDatabase(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Pass string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	descr := data.DBDescription{
		User:     input.Name,
		Pass:     input.Pass,
		FilePath: path.Join(wd, fmt.Sprintf("%s.sqlite", input.Name)),
	}

	v := validator.New()
	if data.ValidateDescription(v, &descr); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.DB.Open(descr)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidAuth):
			app.unauthorizedResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.models.Token.Refresh(); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token := app.models.Token.Get().Plaintext

	err = app.writeJSON(w, http.StatusOK, envelope{"database": input.Name, "authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
