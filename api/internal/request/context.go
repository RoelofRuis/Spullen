package request

import (
	"context"
	"net/http"
)

type contextKey string

const isAuthenticatedKey = contextKey("isAuthenticated")

func ContextSetIsAuthenticated(r *http.Request, isAuthenticated bool) *http.Request {
	ctx := context.WithValue(r.Context(), isAuthenticatedKey, isAuthenticated)
	return r.WithContext(ctx)
}

func ContextGetIsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedKey).(bool)
	if !ok {
		panic("missing isAuthenticated key in request context")
	}
	return isAuthenticated
}
