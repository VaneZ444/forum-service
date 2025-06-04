package handler

import (
	"context"

	"github.com/VaneZ444/forum-service/internal/usecase"
	forumv1 "github.com/VaneZ444/golang-forum-protos/gen/go/forum"
)

type CommentHandler struct {
	forumv1.UnimplementedCommentServiceServer
	uc usecase.CommentUseCase
}

func NewCommentHandler(uc usecase.CommentUseCase) *CommentHandler {
	return &CommentHandler{uc: uc}
}

func (h *CommentHandler) CreateComment(ctx context.Context, req *forumv1.CreateCommentRequest) (*forumv1.CreateCommentResponse, error) {
	id, err := h.uc.CreateComment(ctx, req.PostId, req.AuthorId, req.Content)
	if err != nil {
		return nil, err
	}
	return &forumv1.CreateCommentResponse{CommentId: id}, nil
}

func (h *CommentHandler) GetComment(ctx context.Context, req *forumv1.GetCommentRequest) (*forumv1.GetCommentResponse, error) {
	comment, err := h.uc.GetCommentByID(ctx, req.CommentId)
	if err != nil {
		return nil, err
	}
	return &forumv1.GetCommentResponse{
		Comment: &forumv1.Comment{
			Id:        comment.ID,
			PostId:    comment.PostID,
			Content:   comment.Content,
			AuthorId:  comment.AuthorID,
			CreatedAt: comment.CreatedAt,
		},
	}, nil
}

func (h *CommentHandler) ListCommentsByPost(ctx context.Context, req *forumv1.ListCommentsByPostRequest) (*forumv1.ListCommentsByPostResponse, error) {
	comments, err := h.uc.ListCommentsByPost(ctx, req.PostId)
	if err != nil {
		return nil, err
	}

	var grpcComments []*forumv1.Comment
	for _, c := range comments {
		grpcComments = append(grpcComments, &forumv1.Comment{
			Id:        c.ID,
			PostId:    c.PostID,
			Content:   c.Content,
			AuthorId:  c.AuthorID,
			CreatedAt: c.CreatedAt,
		})
	}

	return &forumv1.ListCommentsByPostResponse{
		Comments: grpcComments,
	}, nil
}
