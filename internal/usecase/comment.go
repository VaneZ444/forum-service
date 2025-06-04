package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

type CommentUseCase interface {
	CreateComment(ctx context.Context, postID int64, authorID int64, content string) (int64, error)
	GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error)
	ListCommentsByPost(ctx context.Context, postID int64) ([]*entity.Comment, error)
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

func (uc *commentUseCase) CreateComment(ctx context.Context, postID int64, authorID int64, content string) (int64, error) {
	_, err := uc.postRepo.GetByID(ctx, postID)
	if err != nil {
		uc.logger.Warn("post not found", slog.Int64("postID", postID), slog.String("err", err.Error()))
		return 0, ErrPostNotFound
	}

	comment := &entity.Comment{
		PostID:    postID,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: time.Now().Unix(),
	}

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

func (uc *commentUseCase) ListCommentsByPost(ctx context.Context, postID int64) ([]*entity.Comment, error) {
	return uc.commentRepo.ListByPostID(ctx, postID, 20, 0)
}
