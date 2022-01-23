package model

import "time"

type Deletion struct {
	DeletedAt   time.Time `json:"deleted_at"`
	Description string    `json:"description"`
}
