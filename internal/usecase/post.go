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
	CreatePost(ctx context.Context, topicID int64, authorID int64, content string) (int64, error)
	GetPostByID(ctx context.Context, id int64) (*entity.Post, error)
	ListPostsByTopic(ctx context.Context, topicID int64) ([]*entity.Post, error)
}

type postUseCase struct {
	postRepo  repository.PostRepository
	topicRepo repository.TopicRepository
	logger    *slog.Logger
}

func NewPostUseCase(postRepo repository.PostRepository, topicRepo repository.TopicRepository, logger *slog.Logger) PostUseCase {
	return &postUseCase{
		postRepo:  postRepo,
		topicRepo: topicRepo,
		logger:    logger,
	}
}

func (uc *postUseCase) CreatePost(ctx context.Context, topicID int64, authorID int64, content string) (int64, error) {
	_, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		uc.logger.Warn("topic not found", slog.Int64("topicID", topicID), slog.String("err", err.Error()))
		return 0, ErrTopicNotFound
	}

	post := &entity.Post{
		TopicID:   topicID,
		AuthorID:  authorID,
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
		uc.logger.Warn("post not found", slog.Int64("id", id), slog.String("err", err.Error()))
		return nil, ErrPostNotFound
	}
	return post, nil
}

func (uc *postUseCase) ListPostsByTopic(ctx context.Context, topicID int64) ([]*entity.Post, error) {
	return uc.postRepo.ListByTopic(ctx, topicID, 20, 0)
}
