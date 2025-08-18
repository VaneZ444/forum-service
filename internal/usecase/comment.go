package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

type CommentUseCase interface {
	CreateComment(ctx context.Context, comment *entity.Comment) (int64, error)
	GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error)
	ListByPost(ctx context.Context, postID int64, limit, offset int) ([]*entity.Comment, int64, error)
}

type commentUseCase struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	logger      *slog.Logger
}

func NewCommentUseCase(commentRepo repository.CommentRepository, postRepo repository.PostRepository, logger *slog.Logger) CommentUseCase {
	return &commentUseCase{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		logger:      logger,
	}
}

func (uc *commentUseCase) CreateComment(ctx context.Context, comment *entity.Comment) (int64, error) {
	_, err := uc.postRepo.GetByID(ctx, comment.PostID)
	if err != nil {
		uc.logger.Warn("post not found", slog.Int64("postID", comment.PostID), slog.String("err", err.Error()))
		return 0, ErrPostNotFound
	}

	comment.CreatedAt = time.Now().Unix()

	id, err := uc.commentRepo.Create(ctx, comment)
	if err != nil {
		uc.logger.Error("failed to create comment", slog.String("err", err.Error()))
		return 0, err
	}

	return id, nil
}

func (uc *commentUseCase) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	comment, err := uc.commentRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Warn("comment not found", slog.Int64("id", id), slog.String("err", err.Error()))
		return nil, ErrCommentNotFound
	}
	return comment, nil
}

func (uc *commentUseCase) ListByPost(ctx context.Context, postID int64, limit, offset int) ([]*entity.Comment, int64, error) {
	if limit <= 0 || limit > 100 {
		return nil, 0, ErrInvalidLimit
	}
	if offset < 0 {
		return nil, 0, ErrInvalidOffset
	}
	// опционально: валидация существования поста
	if _, err := uc.postRepo.GetByID(ctx, postID); err != nil {
		uc.logger.Warn("post not found", slog.Int64("postID", postID))
		return nil, 0, ErrPostNotFound
	}
	return uc.commentRepo.ListByPost(ctx, postID, limit, offset)
}
