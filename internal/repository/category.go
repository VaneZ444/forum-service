package repository

import (
	"context"

	"github.com/VaneZ444/forum-service/internal/entity"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Category, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id int64) error
}
