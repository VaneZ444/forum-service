// internal/usecase/post.go
package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

type PostUseCase interface {
	CreatePost(ctx context.Context, topicID int64, authorID int64, title, content string) (int64, error)
	GetPostByID(ctx context.Context, id int64) (*entity.Post, error)
	ListPostsByTopic(ctx context.Context, topicID int64, limit, offset int) ([]*entity.Post, error)
	ListPosts(ctx context.Context, limit, offset int) ([]*entity.Post, error)
	UpdatePost(ctx context.Context, id int64, title, content string) (*entity.Post, error)
	DeletePost(ctx context.Context, id int64) error
	ListPostsByTag(ctx context.Context, tagID int64, limit, offset int) ([]*entity.Post, error)
}

type postUseCase struct {
	postRepo  repository.PostRepository
	topicRepo repository.TopicRepository
	tagRepo   repository.TagRepository
	logger    *slog.Logger
}

func NewPostUseCase(
	postRepo repository.PostRepository,
	topicRepo repository.TopicRepository,
	tagRepo repository.TagRepository,
	logger *slog.Logger,
) PostUseCase {
	return &postUseCase{
		postRepo:  postRepo,
		topicRepo: topicRepo,
		tagRepo:   tagRepo,
		logger:    logger,
	}
}

func (uc *postUseCase) CreatePost(ctx context.Context, topicID int64, authorID int64, title, content string) (int64, error) {
	_, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		uc.logger.Warn("topic not found", slog.Int64("topicID", topicID))
		return 0, ErrTopicNotFound
	}

	post := &entity.Post{
		TopicID:   topicID,
		AuthorID:  authorID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}

	id, err := uc.postRepo.Create(ctx, post)
	if err != nil {
		uc.logger.Error("failed to create post", slog.String("err", err.Error()))
		return 0, err
	}

	return id, nil
}

func (uc *postUseCase) GetPostByID(ctx context.Context, id int64) (*entity.Post, error) {
	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Warn("post not found", slog.Int64("id", id))
		return nil, ErrPostNotFound
	}
	return post, nil
}

func (uc *postUseCase) ListPostsByTopic(ctx context.Context, topicID int64, limit, offset int) ([]*entity.Post, error) {
	if limit <= 0 || limit > 100 {
		return nil, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, ErrInvalidOffset
	}
	return uc.postRepo.ListByTopic(ctx, topicID, limit, offset)
}

func (uc *postUseCase) ListPosts(ctx context.Context, limit, offset int) ([]*entity.Post, error) {
	if limit <= 0 || limit > 100 {
		return nil, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, ErrInvalidOffset
	}
	return uc.postRepo.List(ctx, limit, offset)
}

func (uc *postUseCase) UpdatePost(ctx context.Context, id int64, title, content string) (*entity.Post, error) {
	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Warn("post not found for update", slog.Int64("id", id))
		return nil, ErrPostNotFound
	}

	post.Title = title
	post.Content = content
	err = uc.postRepo.Update(ctx, post)
	if err != nil {
		uc.logger.Error("failed to update post", slog.String("err", err.Error()))
		return nil, ErrUpdateFailed
	}

	return post, nil
}

func (uc *postUseCase) DeletePost(ctx context.Context, id int64) error {
	err := uc.postRepo.Delete(ctx, id)
	if err != nil {
		uc.logger.Error("failed to delete post", slog.String("err", err.Error()))
		return ErrDeleteFailed
	}
	return nil
}

func (uc *postUseCase) ListPostsByTag(ctx context.Context, tagID int64, limit, offset int) ([]*entity.Post, error) {
	if limit <= 0 || limit > 100 {
		return nil, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, ErrInvalidOffset
	}
	return uc.postRepo.ListByTag(ctx, tagID, limit, offset)
}
