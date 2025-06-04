package handler

import (
	"context"

	"github.com/VaneZ444/forum-service/internal/usecase"
	forumv1 "github.com/VaneZ444/golang-forum-protos/gen/go/forum"
)

type PostHandler struct {
	forumv1.UnimplementedPostServiceServer
	uc usecase.PostUseCase
}

func NewPostHandler(uc usecase.PostUseCase) *PostHandler {
	return &PostHandler{uc: uc}
}

func (h *PostHandler) CreatePost(ctx context.Context, req *forumv1.CreatePostRequest) (*forumv1.CreatePostResponse, error) {
	id, err := h.uc.CreatePost(ctx, req.GetTopicId(), req.GetAuthorId(), req.GetContent())
	if err != nil {
		return nil, err
	}
	return &forumv1.CreatePostResponse{PostId: id}, nil
}

func (h *PostHandler) GetPost(ctx context.Context, req *forumv1.GetPostRequest) (*forumv1.GetPostResponse, error) {
	post, err := h.uc.GetPostByID(ctx, req.GetPostId())
	if err != nil {
		return nil, err
	}
	return &forumv1.GetPostResponse{
		Post: &forumv1.Post{
			Id:        post.ID,
			TopicId:   post.TopicID,
			AuthorId:  post.AuthorID,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
		},
	}, nil
}

func (h *PostHandler) ListPosts(ctx context.Context, req *forumv1.ListPostsRequest) (*forumv1.ListPostsResponse, error) {
	posts, err := h.uc.ListPostsByTopic(ctx, req.GetTopicId())
	if err != nil {
		return nil, err
	}

	var result []*forumv1.Post
	for _, p := range posts {
		result = append(result, &forumv1.Post{
			Id:        p.ID,
			TopicId:   p.TopicID,
			AuthorId:  p.AuthorID,
			Content:   p.Content,
			CreatedAt: p.CreatedAt,
		})
	}

	return &forumv1.ListPostsResponse{Posts: result}, nil
}
