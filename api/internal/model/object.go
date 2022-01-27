package model

import (
	"github.com/roelofruis/spullen/internal/validator"
	"time"
)

func ValidateObject(v *validator.Validator, obj *Object) {
	v.Check(obj.Name != "", "name", "must not be empty")

}

type ObjectID int64

func NewObject(name string, description string) *Object {
	return &Object{
		ID:              ObjectID(0),
		Name:            name,
		Description:     description,
		QuantityChanges: []*QuantityChange{},
		Tags:            []TagID{},
	}
}

type Object struct {
	ID              ObjectID          `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	QuantityChanges []*QuantityChange `json:"quantity_changes"`
	Tags            []TagID           `json:"tags"`
}

func (o *Object) ChangeQuantity(amount int, description string) {
	o.QuantityChanges = append(o.QuantityChanges, &QuantityChange{
		At:          time.Now(),
		Quantity:    amount,
		Description: description,
	})
}

type QuantityChangeID int64

type QuantityChange struct {
	ID          QuantityChangeID `json:"-"`
	At          time.Time        `json:"at"`
	Quantity    int              `json:"quantity"`
	Description string           `json:"description"`
}
