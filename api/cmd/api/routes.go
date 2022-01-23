package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, app.authenticate, app.withQueryParams)

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/db/open", app.handleOpenDatabase)

	router.HandlerFunc(http.MethodGet, "/v1/objects", app.requireAuthentication(app.handleListObjects))
	router.HandlerFunc(http.MethodPost, "/v1/objects", app.requireAuthentication(app.handleAddObject))

	return standardMiddleware.Then(router)
}
