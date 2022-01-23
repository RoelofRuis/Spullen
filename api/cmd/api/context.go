package main

import (
	"context"
	"github.com/roelofruis/spullen/internal/request"
	"net/http"
	"net/url"
)

type contextKey string

const isAuthenticatedKey = contextKey("isAuthenticated")
const queryKey = contextKey("request")

func contextSetIsAuthenticated(r *http.Request, isAuthenticated bool) *http.Request {
	ctx := context.WithValue(r.Context(), isAuthenticatedKey, isAuthenticated)
	return r.WithContext(ctx)
}

func contextGetIsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedKey).(bool)
	if !ok {
		panic("missing isAuthenticated key in request context")
	}
	return isAuthenticated
}

func contextSetQuery(r *http.Request, values url.Values) *http.Request {
	ctx := context.WithValue(r.Context(), queryKey, request.Query{Values: values})
	return r.WithContext(ctx)
}

func contextGetQuery(r *http.Request) request.Query {
	query, ok := r.Context().Value(queryKey).(request.Query)
	if !ok {
		panic("missing request key in request context")
	}
	return query
}
