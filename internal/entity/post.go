package entity

type Post struct {
	ID        int64
	TopicID   int64
	Title     string
	Content   string
	AuthorID  int64
	CreatedAt int64
}
