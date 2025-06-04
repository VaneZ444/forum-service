package entity

type Comment struct {
	ID        int64
	PostID    int64
	Content   string
	AuthorID  int64
	CreatedAt int64
}
