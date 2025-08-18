package repository

import (
	"context"

	"github.com/VaneZ444/forum-service/internal/entity"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *entity.Comment) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Comment, error)
	ListByPost(ctx context.Context, postID int64, limit, offset int) ([]*entity.Comment, int64, error)
}
