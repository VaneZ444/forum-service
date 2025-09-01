package entity

import "time"

type Comment struct {
	ID        int64
	PostID    int64
	Content   string
	AuthorID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
