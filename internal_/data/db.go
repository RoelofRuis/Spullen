package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/roelofruis/spullen/internal_/validator"
	"os"
	"regexp"
	"sync"

	_ "github.com/CovenantSQL/go-sqlite3-encrypt"
)

var (
	FileRX = regexp.MustCompile("^[a-zA-Z0-9_]*$")
)

var (
	ModeCreate      = "create"
	ModeOpen        = "open"
	ErrNoDataSource = errors.New("no data source opened")
)

type DBDescription struct {
	Name     string
	Key      string
	Mode     string
	FilePath string
}

func ValidateDescription(v *validator.Validator, descr *DBDescription) {
	v.Check(descr.Name != "", "name", "must not be empty")
	v.Check(validator.Matches(descr.Name, FileRX), "name", "can only contain alphanumeric characters and underscore")
	v.Check(descr.Key != "", "key", "must not be empty")
	v.Check(validator.In(descr.Mode, ModeCreate, ModeOpen), "mode", "must be one of 'create' or 'open'")

	if descr.Mode == ModeOpen {
		if _, err := os.Stat(descr.FilePath); errors.Is(err, os.ErrNotExist) {
			v.AddError("name", "does not exist")
		}
	}
}

type DBProxy struct {
	db   *sql.DB
	lock sync.RWMutex
}

func NewDBProxy() *DBProxy {
	return &DBProxy{
		db:   nil,
		lock: sync.RWMutex{},
	}
}

func (db *DBProxy) Open(descr DBDescription) error {
	conn, err := sql.Open(
		"sqlite3",
		fmt.Sprintf(
			"file:%s?_auth&_auth_user=%s&_auth_pass=%s",
			descr.FilePath,
			"admin",
			"admin",
		),
	)
	if err != nil {
		return err
	}

	keyHash := fmt.Sprintf("%x", sha256.Sum256([]byte(descr.Key)))

	_, err = conn.Exec(fmt.Sprintf("PRAGMA key='%s'", keyHash))
	if err != nil {
		return err
	}

	_, err = conn.Exec(`CREATE TABLE test(x integer PRIMARY KEY)`)
	if err != nil {
		return err
	}

	db.db = conn
	return nil
}

func (db *DBProxy) Close() (err error) {
	if db.db != nil {
		err = db.db.Close()
		db.db = nil
	}
	return
}

func (db *DBProxy) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if db.db == nil {
		return nil, ErrNoDataSource
	}

	return db.db.QueryContext(ctx, query, args)
}
