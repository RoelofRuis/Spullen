package data

import (
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"github.com/roelofruis/spullen/internal_/migration"
	"github.com/roelofruis/spullen/internal_/validator"
	"regexp"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	FileRX = regexp.MustCompile("^[a-zA-Z0-9_]*$")
)

var (
	ErrInvalidAuth = errors.New("invalid authorization")
	ErrNoDataSource = errors.New("no data source opened")
)

type DBDescription struct {
	User     string
	Pass     string
	FilePath string
}

func ValidateDescription(v *validator.Validator, descr *DBDescription) {
	v.Check(descr.User != "", "user", "must not be empty")
	v.Check(validator.Matches(descr.User, FileRX), "user", "can only contain alphanumeric characters and underscore")
	v.Check(descr.Pass != "", "pass", "must not be empty")
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
	passHash := md5.Sum([]byte(descr.Pass))

	conn, err := sql.Open(
		"sqlite3",
		fmt.Sprintf(
			"file:%s?_auth&_auth_user=%s&_auth_pass=%s&_auth_crypt=sha256",
			descr.FilePath,
			descr.User,
			passHash,
		),
	)
	if err != nil {
		return err
	}

	migrator, err := migration.Init(conn)
	if err != nil {
		if strings.Contains(err.Error(), "SQLITE_AUTH: Unauthorized") {
			return ErrInvalidAuth
		}
		return err
	}

	if err := migrator.Up(); err != nil {
		return err
	}

	db.lock.Lock()
	db.db = conn
	db.lock.Unlock()
	return nil
}

func (db *DBProxy) Close() (err error) {
	if db.db != nil {
		db.lock.Lock()
		err = db.db.Close()
		db.db = nil
		db.lock.Unlock()
	}
	return
}

func (db *DBProxy) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return nil, ErrNoDataSource
	}

	return db.db.QueryContext(ctx, query, args...)
}
