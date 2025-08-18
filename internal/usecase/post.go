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
	CreatePost(ctx context.Context, post *entity.Post) (int64, error)
	GetPostByID(ctx context.Context, id int64) (*entity.Post, error)
	ListByTopic(ctx context.Context, topicID int64, limit, offset int) ([]*entity.Post, error)
	List(ctx context.Context, topicID, tagID int64, limit, offset int) ([]*entity.Post, int64, error)
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

func (uc *postUseCase) CreatePost(ctx context.Context, post *entity.Post) (int64, error) {
	_, err := uc.topicRepo.GetByID(ctx, post.TopicID)
	if err != nil {
		uc.logger.Warn("topic not found", slog.Int64("topicID", post.TopicID))
		return 0, ErrTopicNotFound
	}

	post.CreatedAt = time.Now().Unix()

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

func (uc *postUseCase) ListByTopic(ctx context.Context, topicID int64, limit, offset int) ([]*entity.Post, error) {
	if limit <= 0 || limit > 100 {
		return nil, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, ErrInvalidOffset
	}
	return uc.postRepo.ListByTopic(ctx, topicID, limit, offset)
}

func (uc *postUseCase) List(ctx context.Context, topicID, tagID int64, limit, offset int) ([]*entity.Post, int64, error) {
	if limit <= 0 || limit > 100 {
		return nil, 0, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, 0, ErrInvalidOffset
	}

	posts, total, err := uc.postRepo.List(ctx, topicID, tagID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
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
