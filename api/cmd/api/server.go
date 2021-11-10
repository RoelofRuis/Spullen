package main

import (
	"errors"
	"net/http"
	"time"
)

func (app *application) serve() error {
	address := "localhost:8080"

	srv := &http.Server{
		Addr:         address,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	app.logger.PrintInfo("Starting API", map[string]string{"address": address})
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
