package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
	"github.com/gosimple/slug"
)

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id int64) (*entity.Category, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Category, int64, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryUseCase struct {
	categoryRepo repository.CategoryRepository
	logger       *slog.Logger
}

func NewCategoryUseCase(categoryRepo repository.CategoryRepository, logger *slog.Logger) CategoryUseCase {
	return &categoryUseCase{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *categoryUseCase) CreateCategory(ctx context.Context, category *entity.Category) error {
	// Generate slug if not provided
	if category.Slug == "" {
		category.Slug = slug.Make(category.Title)
	}

	// Check if slug exists
	existing, err := uc.categoryRepo.GetBySlug(ctx, category.Slug)
	if err == nil && existing != nil {
		uc.logger.Warn("category already exists", slog.String("slug", category.Slug))
		return ErrCategoryAlreadyExists
	}

	// Set timestamps
	now := time.Now().UnixMilli()
	category.CreatedAt = now
	category.UpdatedAt = now

	return uc.categoryRepo.Create(ctx, category)
}

func (uc *categoryUseCase) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return category, nil
}

func (uc *categoryUseCase) List(ctx context.Context, limit, offset int) ([]*entity.Category, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	categories, err := uc.categoryRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.categoryRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (uc *categoryUseCase) UpdateCategory(ctx context.Context, category *entity.Category) error {
	existing, err := uc.categoryRepo.GetByID(ctx, category.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			return ErrCategoryNotFound
		}
		return err
	}

	// Preserve created_at
	category.CreatedAt = existing.CreatedAt
	category.UpdatedAt = time.Now().UnixMilli()

	return uc.categoryRepo.Update(ctx, category)
}

func (uc *categoryUseCase) DeleteCategory(ctx context.Context, id int64) error {
	return uc.categoryRepo.Delete(ctx, id)
}
