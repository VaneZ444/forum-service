package entity

import "time"

type Post struct {
	ID            int64
	TopicID       int64
	AuthorID      int64
	Title         string
	Content       string
	Images        []string
	Tags          []Tag
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Status        Status
	ViewsCount    int64
	CommentsCount int64
	LikesCount    int64
}
