package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	db *sql.DB
}

func NewSqliteStorage() (*SqliteStorage, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	return &SqliteStorage{db}, nil
}

func (ol *SqliteStorage) Get(id string) *Object {
	// TODO: implement

	return nil
}

func (ol *SqliteStorage) GetAll() *map[string]*Object {
	// TODO: implement

	return nil
}

func (ol *SqliteStorage) PutObject(o *Object) error {
	// TODO: implement

	return nil
}

func (ol *SqliteStorage) RemoveObject(id string) error {
	// TODO: implement

	return nil
}
