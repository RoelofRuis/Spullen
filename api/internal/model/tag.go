package model

import "github.com/roelofruis/spullen/internal/validator"

func ValidateTag(v *validator.Validator, tag *Tag) {
	v.Check(tag.Name != "", "name", "must not be empty")
}

type TagID int64

type Tag struct {
	ID          TagID  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsSystemTag bool   `json:"is_system_tag"`
}
