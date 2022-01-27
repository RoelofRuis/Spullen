package model

import (
	"github.com/roelofruis/spullen/internal/validator"
	"time"
)

func ValidateObject(v *validator.Validator, obj *Object) {
	v.Check(obj.Name != "", "name", "must not be empty")
	v.Check(obj.Quantity > 0, "quantity", "quantity must be a positive integer")
}

type ObjectID int64

type Object struct {
	ID       ObjectID  `json:"id"`
	Added    time.Time `json:"added"`
	Name     string    `json:"name"`
	Quantity int       `json:"quantity"`
	Deletion *Deletion `json:"deletion"`
	Tag      []TagID   `json:"tags"`
}

type Deletion struct {
	DeletedAt   time.Time `json:"deleted_at"`
	Description string    `json:"description"`
}

func (o *Object) Delete(at time.Time, description string) {
	o.Deletion = &Deletion{
		DeletedAt:   at,
		Description: description,
	}
}
