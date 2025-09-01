package entity

import "time"

type Category struct {
	ID          int64
	Title       string
	Slug        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
