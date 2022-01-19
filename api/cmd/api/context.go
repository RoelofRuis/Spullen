package main

import (
	"context"
	"net/http"
)

type contextKey string

const isAuthenticatedKey = contextKey("isAuthenticated")

func (app *application) contextSetIsAuthenticated(r *http.Request, isAuthenticated bool) *http.Request {
	ctx := context.WithValue(r.Context(), isAuthenticatedKey, isAuthenticated)
	return r.WithContext(ctx)
}

func (app *application) contextGetIsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedKey).(bool)
	if !ok {
		panic("missing isAuthenticated key in request context")
	}
	return isAuthenticated
}
