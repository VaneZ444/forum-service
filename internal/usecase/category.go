package usecase

import (
	"context"
	"log/slog"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, title, slug string) (int64, error)
	GetCategoryByID(ctx context.Context, id int64) (*entity.Category, error)
	ListCategories(ctx context.Context) ([]*entity.Category, error)
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

func (uc *categoryUseCase) GetCategoryByID(ctx context.Context, id int64) (*entity.Category, error) {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Warn("category not found", slog.Int64("id", id), slog.String("err", err.Error()))
		return nil, ErrCategoryNotFound
	}

	return category, nil
}

func (uc *categoryUseCase) ListCategories(ctx context.Context) ([]*entity.Category, error) {
	return uc.categoryRepo.List(ctx)
}
