package repository

import (
	"context"

	"github.com/VaneZ444/forum-service/internal/entity"
)

type CategoryRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.Category, error)
	List(ctx context.Context) ([]*entity.Category, error)
	Create(ctx context.Context, category *entity.Category) (int64, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
}
