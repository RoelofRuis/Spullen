package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/CovenantSQL/go-sqlite3-encrypt"
)

var (
	ErrNoDataSource = errors.New("no data source opened")
)

type DBProxy struct {
	DB *sql.DB
}

func (db *DBProxy) Open(dbName string, key string) error {
	fileName := fmt.Sprintf("%s.sqlite", dbName)
	conn, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", fileName))
	if err != nil {
		return err
	}

	keyHash := fmt.Sprintf("%x", sha256.Sum256([]byte(key)))

	_, err = conn.Exec(fmt.Sprintf("PRAGMA key='%s'", keyHash))
	if err != nil {
		return err
	}

	db.DB = conn
	return nil
}

func (db *DBProxy) Close() (err error) {
	if db.DB != nil {
		err = db.DB.Close()
		db.DB = nil
	}
	return
}

func (db *DBProxy) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if db.DB == nil {
		return nil, ErrNoDataSource
	}

	return db.DB.QueryContext(ctx, query, args)
}
