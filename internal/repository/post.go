package repository

import (
	"context"

	"github.com/VaneZ444/forum-service/internal/entity"
)

type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Post, error)
	ListByTopic(ctx context.Context, topicID int64, limit int, offset int) ([]*entity.Post, error)
	List(ctx context.Context, topicID, tagID int64, limit, offset int) ([]*entity.Post, int64, error)
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, id int64) error
	ListByTag(ctx context.Context, tagID int64, limit, offset int) ([]*entity.Post, error)
}
