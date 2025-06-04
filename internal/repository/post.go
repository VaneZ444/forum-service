package repository

import (
	"context"

	"github.com/VaneZ444/forum-service/internal/entity"
)

type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Post, error)
	ListByTopic(ctx context.Context, topicID int64, limit int, offset int) ([]*entity.Post, error)
}
