package entity

import "time"

type Topic struct {
	ID           int64
	Title        string
	AuthorID     int64
	CategoryID   int64
	CreatedAt    time.Time
	Status       Status // ACTIVE, DELETED, etc.
	PostsCount   int64
	ViewsCount   int64
	LastActivity time.Time
}
