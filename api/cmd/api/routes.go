package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest)

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/db/open", app.handleOpenDatabase)

	router.HandlerFunc(http.MethodGet, "/v1/objects", app.handleListObjects)
	router.HandlerFunc(http.MethodPost, "/v1/objects", app.handleAddObject)

	return standardMiddleware.Then(router)
}