package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

var (
	ErrTagNotFound = errors.New("tag not found")
)

type TagUseCase interface {
	CreateTag(ctx context.Context, name string) (int64, error)
	GetTagByID(ctx context.Context, id int64) (*entity.Tag, error)
	ListTags(ctx context.Context) ([]*entity.Tag, error)
	ListTagsByPostID(ctx context.Context, postID int64) ([]*entity.Tag, error)
}

type tagUseCase struct {
	tagRepo repository.TagRepository
	logger  *slog.Logger
}

func NewTagUseCase(tagRepo repository.TagRepository, logger *slog.Logger) TagUseCase {
	return &tagUseCase{
		tagRepo: tagRepo,
		logger:  logger,
	}
}

func (uc *tagUseCase) CreateTag(ctx context.Context, name string) (int64, error) {
	tag := &entity.Tag{Name: name}
	id, err := uc.tagRepo.Create(ctx, tag)
	if err != nil {
		uc.logger.Error("failed to create tag", slog.String("err", err.Error()))
		return 0, err
	}
	return id, nil
}

func (uc *tagUseCase) GetTagByID(ctx context.Context, id int64) (*entity.Tag, error) {
	tag, err := uc.tagRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Warn("tag not found", slog.Int64("id", id), slog.String("err", err.Error()))
		return nil, ErrTagNotFound
	}
	return tag, nil
}

func (uc *tagUseCase) ListTags(ctx context.Context) ([]*entity.Tag, error) {
	return uc.tagRepo.List(ctx)
}

func (uc *tagUseCase) ListTagsByPostID(ctx context.Context, postID int64) ([]*entity.Tag, error) {
	return uc.tagRepo.ListByPostID(ctx, postID)
}
