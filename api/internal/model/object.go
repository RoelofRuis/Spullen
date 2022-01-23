package model

import "time"

type Object struct {
	ID       int64     `json:"id"`
	Added    time.Time `json:"added"`
	Name     string    `json:"name"`
	Quantity int       `json:"quantity"`
	Deletion *Deletion `json:"deletion"`
	Tag      []*Tag    `json:"tags"`
}

func (o *Object) Delete(at time.Time, description string) {
	o.Deletion = &Deletion{
		DeletedAt:   at,
		Description: description,
	}
}
