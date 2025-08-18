package usecase

import (
	"context"
	"log/slog"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, title, description string) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Category, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Category, error)
	UpdateCategory(ctx context.Context, id int64, title, description string) (*entity.Category, error)
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

func (uc *categoryUseCase) CreateCategory(ctx context.Context, title, slug string) (int64, error) {
	existing, err := uc.categoryRepo.GetBySlug(ctx, slug)
	if err == nil && existing != nil {
		uc.logger.Warn("category already exists", slog.String("slug", slug))
		return 0, ErrCategoryAlreadyExists
	}

	category := &entity.Category{
		Title: title,
		Slug:  slug,
	}

	id, err := uc.categoryRepo.Create(ctx, category)
	if err != nil {
		uc.logger.Error("failed to create category", slog.String("err", err.Error()))
		return 0, err
	}

	return id, nil
}

func (uc *categoryUseCase) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (uc *categoryUseCase) List(ctx context.Context, limit, offset int) ([]*entity.Category, error) {
	if limit <= 0 || limit > 100 {
		return nil, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, ErrInvalidOffset
	}
	return uc.categoryRepo.List(ctx, limit, offset)
}

func (uc *categoryUseCase) UpdateCategory(ctx context.Context, id int64, title, description string) (*entity.Category, error) {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Warn("category not found for update", slog.Int64("id", id))
		return nil, ErrCategoryNotFound
	}

	category.Title = title
	category.Description = description

	err = uc.categoryRepo.Update(ctx, category)
	if err != nil {
		uc.logger.Error("failed to update category", slog.String("err", err.Error()))
		return nil, ErrUpdateFailed
	}

	return category, nil
}

func (uc *categoryUseCase) DeleteCategory(ctx context.Context, id int64) error {
	err := uc.categoryRepo.Delete(ctx, id)
	if err != nil {
		uc.logger.Error("failed to delete category", slog.String("err", err.Error()))
		return ErrDeleteFailed
	}
	return nil
}
