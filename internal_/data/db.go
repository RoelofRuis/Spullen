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

func (db DBProxy) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if db.DB == nil {
		return nil, ErrNoDataSource
	}

	return db.DB.QueryContext(ctx, query, args)
}
