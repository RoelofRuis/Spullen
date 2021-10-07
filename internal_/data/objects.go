package data

import (
	"database/sql"
	"time"
)

type Object struct {
	ID       int64
	Added    time.Time
	Name     string
	Quantity int
}

type ObjectRepository struct {
	DB *sql.DB
}

func (r ObjectRepository) GetAll() []*Object {
	return nil
}
