package data

import (
	"time"
)

type Object struct {
	ID       int64
	Added    time.Time
	Name     string
	Quantity int
}

type ObjectModel struct {
	DB DBProxy
}

func (r ObjectModel) GetAll() ([]*Object, error) {
	return nil, nil
}
