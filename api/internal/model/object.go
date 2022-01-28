package model

import (
	"github.com/roelofruis/spullen/internal/validator"
	"time"
)

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

func (o *Object) Validate(v *validator.Validator) {
	v.Check(o.Name != "", "name", "must not be empty")
	v.Check(o.GetTotalQuantity() >= 0, "total_quantity", "cannot be negative")
}

func (o *Object) ChangeQuantity(amount int, description string) {
	o.QuantityChanges = append(o.QuantityChanges, &QuantityChange{
		At:          time.Now(),
		Quantity:    amount,
		Description: description,
	})
}

func (o *Object) AttachTag(t *Tag) {
	for _, id := range o.Tags {
		if id == t.ID {
			return
		}
	}
	o.Tags = append(o.Tags, t.ID)
}

func (o *Object) RemoveTag(t *Tag) {
	for i, id := range o.Tags {
		if id == t.ID {
			o.Tags[i] = o.Tags[len(o.Tags)-1]
			o.Tags = o.Tags[:len(o.Tags)-1]
			return
		}
	}
}

func (o *Object) GetTotalQuantity() int {
	sum := 0
	for _, c := range o.QuantityChanges {
		sum += c.Quantity
	}
	return sum
}

type QuantityChangeID int64

type QuantityChange struct {
	ID          QuantityChangeID `json:"-"`
	At          time.Time        `json:"at"`
	Quantity    int              `json:"quantity"`
	Description string           `json:"description"`
}
