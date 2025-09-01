package repository

import (
	"context"
	"errors"

	"github.com/VaneZ444/forum-service/internal/entity"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) (*entity.Category, error)
	GetByID(ctx context.Context, id int64) (*entity.Category, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Category, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, category *entity.Category) (*entity.Category, error)
	Delete(ctx context.Context, id int64) error
}
