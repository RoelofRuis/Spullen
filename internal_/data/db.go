package data

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNoDataSource = errors.New("no data source opened")
)

type DBProxy struct {
	DB *sql.DB
}

func (db DBProxy) Open(dns string) error {
	return nil // TODO: implement
}

func (db DBProxy) Close() error {
	return nil // TODO: implement
}

func (db DBProxy) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if db.DB == nil {
		return nil, ErrNoDataSource
	}

	return db.DB.QueryContext(ctx, query, args)
}
