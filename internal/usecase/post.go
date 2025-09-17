// internal/usecase/post.go
package usecase

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
	forumv1 "github.com/VaneZ444/golang-forum-protos/gen/go/forum"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostUseCase interface {
	CreatePost(ctx context.Context, post *entity.Post) (int64, error)
	GetPostByID(ctx context.Context, id int64) (*entity.Post, error)
	ListByTopic(ctx context.Context, topicID int64, limit, offset int) ([]*entity.Post, int64, error)
	List(ctx context.Context, topicID, tagID int64, limit, offset int) ([]*entity.Post, int64, error)
	UpdatePost(ctx context.Context, req *forumv1.UpdatePostRequest) (*entity.Post, error)
	DeletePost(ctx context.Context, id int64) error
	ListPostsByTag(ctx context.Context, tagID int64, limit, offset int) ([]*entity.Post, int64, error)
	AddView(ctx context.Context, postID, userID int64) error
	SearchPosts(ctx context.Context, query string, limit, offset int) ([]*entity.Post, int64, error)
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

	post.CreatedAt = time.Now().UTC()

	id, err := uc.postRepo.Create(ctx, post)
	if err != nil {
		uc.logger.Error("failed to create post", slog.String("err", err.Error()))
		return 0, err
	}

	return id, nil
}
func (uc *postUseCase) AddView(ctx context.Context, postID, userID int64) error {
	return uc.postRepo.AddView(ctx, postID, userID)
}
func (uc *postUseCase) GetPostByID(ctx context.Context, id int64) (*entity.Post, error) {
	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Warn("post not found", slog.Int64("id", id))
		return nil, ErrPostNotFound
	}
	return post, nil
}

func (uc *postUseCase) ListByTopic(ctx context.Context, topicID int64, limit, offset int) ([]*entity.Post, int64, error) {
	if limit <= 0 || limit > 100 {
		return nil, 0, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, 0, ErrInvalidOffset
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
func (uc *postUseCase) SearchPosts(ctx context.Context, query string, limit, offset int) ([]*entity.Post, int64, error) {
	if limit <= 0 || limit > 100 {
		return nil, 0, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, 0, ErrInvalidOffset
	}
	return uc.postRepo.Search(ctx, query, limit, offset)
}

func (uc *postUseCase) UpdatePost(ctx context.Context, req *forumv1.UpdatePostRequest) (*entity.Post, error) {
	post, err := uc.postRepo.GetByID(ctx, req.GetId())
	if err != nil {
		uc.logger.Warn("post not found", slog.Int64("id", req.GetId()))
		return nil, ErrPostNotFound
	}

	// Мержим изменения
	if req.Title != nil {
		post.Title = strings.TrimSpace(req.GetTitle())
		if post.Title == "" {
			return nil, status.Error(codes.InvalidArgument, "title cannot be empty")
		}
	}
	if req.Content != nil {
		post.Content = req.GetContent()
	}
	if req.Images != nil {
		post.Images = req.GetImages()
	}
	if req.TagIds != nil {
		post.Tags = make([]entity.Tag, len(req.GetTagIds()))
		for i, id := range req.GetTagIds() {
			post.Tags[i] = entity.Tag{ID: id}
		}
	}

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

func (uc *postUseCase) ListPostsByTag(ctx context.Context, tagID int64, limit, offset int) ([]*entity.Post, int64, error) {
	if limit <= 0 || limit > 100 {
		return nil, 0, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, 0, ErrInvalidOffset
	}

	posts, total, err := uc.postRepo.ListByTag(ctx, tagID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
