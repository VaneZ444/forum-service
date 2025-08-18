package usecase

import (
	"context"
	"log/slog"
	"strings"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

type TagUseCase interface {
	CreateTag(ctx context.Context, tag *entity.Tag) error
	GetTagByID(ctx context.Context, id int64) (*entity.Tag, error)
	GetTagBySlug(ctx context.Context, slug string) (*entity.Tag, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Tag, int64, error)
	ListTagsByPostID(ctx context.Context, postID int64) ([]*entity.Tag, error)
	AddTagToPost(ctx context.Context, postID, tagID int64) error
	RemoveTagFromPost(ctx context.Context, postID, tagID int64) error
}

type tagUseCase struct {
	tagRepo  repository.TagRepository
	postRepo repository.PostRepository
	logger   *slog.Logger
}

func NewTagUseCase(
	tagRepo repository.TagRepository,
	postRepo repository.PostRepository,
	logger *slog.Logger,
) TagUseCase {
	return &tagUseCase{
		tagRepo:  tagRepo,
		postRepo: postRepo,
		logger:   logger,
	}
}

func (uc *tagUseCase) CreateTag(ctx context.Context, tag *entity.Tag) error {
	if tag.Slug == "" {
		tag.Slug = strings.ToLower(strings.ReplaceAll(tag.Name, " ", "-"))
	}

	id, err := uc.tagRepo.Create(ctx, tag)
	if err != nil {
		uc.logger.Error("failed to create tag", slog.String("err", err.Error()))
		return err
	}

	tag.ID = id
	return nil
}

func (uc *tagUseCase) GetTagByID(ctx context.Context, id int64) (*entity.Tag, error) {
	tag, err := uc.tagRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Warn("tag not found", slog.Int64("id", id))
		return nil, ErrTagNotFound
	}
	return tag, nil
}
func (uc *tagUseCase) GetTagBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	tag, err := uc.tagRepo.GetBySlug(ctx, slug)
	if err != nil {
		uc.logger.Warn("tag not found", slog.String("slug", slug))
		return nil, ErrTagNotFound
	}
	return tag, nil
}
func (uc *tagUseCase) List(ctx context.Context, limit, offset int) ([]*entity.Tag, int64, error) {
	if limit <= 0 || limit > 100 {
		return nil, 0, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, 0, ErrInvalidOffset
	}
	return uc.tagRepo.List(ctx, limit, offset)
}

func (uc *tagUseCase) ListTagsByPostID(ctx context.Context, postID int64) ([]*entity.Tag, error) {
	_, err := uc.postRepo.GetByID(ctx, postID)
	if err != nil {
		uc.logger.Warn("post not found", slog.Int64("postID", postID))
		return nil, ErrPostNotFound
	}
	return uc.tagRepo.ListByPostID(ctx, postID)
}

func (uc *tagUseCase) AddTagToPost(ctx context.Context, postID, tagID int64) error {
	_, err := uc.postRepo.GetByID(ctx, postID)
	if err != nil {
		uc.logger.Warn("post not found", slog.Int64("postID", postID))
		return ErrPostNotFound
	}

	_, err = uc.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		uc.logger.Warn("tag not found", slog.Int64("tagID", tagID))
		return ErrTagNotFound
	}

	err = uc.tagRepo.AddToPost(ctx, postID, tagID)
	if err != nil {
		uc.logger.Error("failed to add tag to post", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (uc *tagUseCase) RemoveTagFromPost(ctx context.Context, postID, tagID int64) error {
	err := uc.tagRepo.RemoveFromPost(ctx, postID, tagID)
	if err != nil {
		uc.logger.Error("failed to remove tag from post", slog.String("err", err.Error()))
		return err
	}
	return nil
}
