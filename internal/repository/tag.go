package repository

import (
	"context"

	"github.com/VaneZ444/forum-service/internal/entity"
)

type TagRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.Tag, error)
	ListByIDs(ctx context.Context, ids []int64) ([]*entity.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Tag, error)
	Create(ctx context.Context, tag *entity.Tag) (int64, error)
	ListAll(ctx context.Context) ([]*entity.Tag, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Tag, error)
	ListByPostID(ctx context.Context, postID int64) ([]*entity.Tag, error)
	AddToPost(ctx context.Context, postID int64, tagID int64) error
	RemoveFromPost(ctx context.Context, postID int64, tagID int64) error
}
