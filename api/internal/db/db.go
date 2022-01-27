package db

import (
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/roelofruis/spullen/internal/migration"
	"github.com/roelofruis/spullen/internal/validator"
	"regexp"
	"strings"
	"sync"
	"time"
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

type Proxy struct {
	db   *sql.DB
	lock sync.RWMutex
}

func NewProxy() *Proxy {
	return &Proxy{
		db:   nil,
		lock: sync.RWMutex{},
	}
}

func (p *Proxy) Open(descr DBDescription) error {
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

	p.lock.Lock()
	p.db = conn
	p.lock.Unlock()
	return nil
}

func (p *Proxy) Close() (err error) {
	if p.db != nil {
		p.lock.Lock()
		err = p.db.Close()
		p.db = nil
		p.lock.Unlock()
	}
	return
}

func (p *Proxy) Exec(i *InsertStatement) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.ExecContext(ctx, i.query(), i.args()...)
}

func (p *Proxy) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if p.db == nil {
		return nil, ErrNoDataSource
	}

	return p.db.ExecContext(ctx, query, args...)
}

func (p *Proxy) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if p.db == nil {
		return nil, ErrNoDataSource
	}

	return p.db.QueryContext(ctx, query, args...)
}
