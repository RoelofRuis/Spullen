package model

import (
	"github.com/roelofruis/spullen/internal/validator"
)

func ValidateObject(v *validator.Validator, obj *Object) {
	v.Check(obj.Name != "", "name", "must not be empty")
}

type ObjectID int64

type Object struct {
	ID          ObjectID  `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tag         []TagID   `json:"tags"`
}
