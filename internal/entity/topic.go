package entity

type Topic struct {
	ID           int64
	Title        string
	AuthorID     int64
	CategoryID   int64
	CreatedAt    int64
	Status       Status // ACTIVE, DELETED, etc.
	PostsCount   int64
	ViewsCount   int64
	LastActivity int64
}
