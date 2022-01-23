package model

type Tag struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsSystemTag bool   `json:"is_system_tag"`
}
