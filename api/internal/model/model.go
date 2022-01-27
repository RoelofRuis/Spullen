package model

import (
	"github.com/roelofruis/spullen/internal/db"
)

type Model struct {
	DB      *db.Proxy
	Token   *Token
	Objects *ObjectRepository
	Tags    *TagRepository
}

func NewModel(db *db.Proxy) Model {
	return Model{
		DB:      db,
		Token:   &Token{},
		Objects: &ObjectRepository{DB: db},
		Tags:    &TagRepository{DB: db},
	}
}
