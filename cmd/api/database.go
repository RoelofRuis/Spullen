package main

import (
	"fmt"
	"net/http"
)

func (app *application) handleNewDatabase(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create new database")
}
