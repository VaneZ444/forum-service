package entity

type Post struct {
	ID            int64
	TopicID       int64
	AuthorID      int64
	Title         string
	Content       string
	Images        []string
	Tags          []Tag
	CreatedAt     int64
	UpdatedAt     int64
	Status        Status
	ViewsCount    int64
	CommentsCount int64
	LikesCount    int64
}
