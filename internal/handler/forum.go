package handler

import (
	"context"
	"log/slog"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/usecase"
	forumv1 "github.com/VaneZ444/golang-forum-protos/gen/go/forum"
)

type ForumHandler struct {
	forumv1.UnimplementedForumServer
	categoryUC usecase.CategoryUseCase
	topicUC    usecase.TopicUseCase
	postUC     usecase.PostUseCase
	commentUC  usecase.CommentUseCase
	tagUC      usecase.TagUseCase
	logger     *slog.Logger
}

func NewForumHandler(
	categoryUC usecase.CategoryUseCase,
	topicUC usecase.TopicUseCase,
	postUC usecase.PostUseCase,
	commentUC usecase.CommentUseCase,
	tagUC usecase.TagUseCase,
	logger *slog.Logger,
) *ForumHandler {
	return &ForumHandler{
		categoryUC: categoryUC,
		topicUC:    topicUC,
		postUC:     postUC,
		commentUC:  commentUC,
		tagUC:      tagUC,
		logger:     logger,
	}
}

func (h *ForumHandler) CreateCategory(ctx context.Context, req *forumv1.CreateCategoryRequest) (*forumv1.CreateCategoryResponse, error) {
	h.logger.Info("creating category", "title", req.Title)
	id, err := h.categoryUC.CreateCategory(ctx, req.Title, req.Description)
	if err != nil {
		h.logger.Error("failed to create category", "error", err)
		return nil, err
	}
	return &forumv1.CreateCategoryResponse{CategoryId: id}, nil
}

func (h *ForumHandler) GetCategoryByID(ctx context.Context, req *forumv1.GetCategoryByIDRequest) (*forumv1.GetCategoryByIDResponse, error) {
	category, err := h.categoryUC.GetByID(ctx, req.CategoryId)
	if err != nil {
		h.logger.Error("failed to get category", "error", err)
		return nil, err
	}
	return &forumv1.GetCategoryByIDResponse{Category: toProtoCategory(category)}, nil
}

func (h *ForumHandler) ListCategories(ctx context.Context, _ *forumv1.Empty) (*forumv1.ListCategoriesResponse, error) {
	limit := 100
	offset := 0
	categories, err := h.categoryUC.List(ctx, limit, offset)
	if err != nil {
		h.logger.Error("failed to list categories", "error", err)
		return nil, err
	}

	protoCategories := make([]*forumv1.Category, len(categories))
	for i, c := range categories {
		protoCategories[i] = toProtoCategory(c)
	}

	return &forumv1.ListCategoriesResponse{Categories: protoCategories}, nil
}

func (h *ForumHandler) CreateTopic(ctx context.Context, req *forumv1.CreateTopicRequest) (*forumv1.CreateTopicResponse, error) {
	h.logger.Info("creating topic", "title", req.Title)
	id, err := h.topicUC.CreateTopic(ctx, req.Title, req.CategoryId, req.AuthorId)
	if err != nil {
		h.logger.Error("failed to create topic", "error", err)
		return nil, err
	}
	return &forumv1.CreateTopicResponse{TopicId: id}, nil
}

func (h *ForumHandler) GetTopicByID(ctx context.Context, req *forumv1.GetTopicByIDRequest) (*forumv1.GetTopicByIDResponse, error) {
	topic, err := h.topicUC.GetByID(ctx, req.TopicId)
	if err != nil {
		h.logger.Error("failed to get topic", "error", err)
		return nil, err
	}
	return &forumv1.GetTopicByIDResponse{Topic: toProtoTopic(topic)}, nil
}

func (h *ForumHandler) ListTopicsByCategory(ctx context.Context, req *forumv1.ListTopicsByCategoryRequest) (*forumv1.ListTopicsByCategoryResponse, error) {
	limit := int(req.Limit)
	offset := int(req.Offset)
	topics, err := h.topicUC.ListByCategory(ctx, req.CategoryId, limit, offset)
	if err != nil {
		h.logger.Error("failed to list topics", "error", err)
		return nil, err
	}

	protoTopics := make([]*forumv1.Topic, len(topics))
	for i, t := range topics {
		protoTopics[i] = toProtoTopic(t)
	}

	return &forumv1.ListTopicsByCategoryResponse{Topics: protoTopics}, nil
}

func (h *ForumHandler) CreatePost(ctx context.Context, req *forumv1.CreatePostRequest) (*forumv1.CreatePostResponse, error) {
	h.logger.Info("creating post", "topic_id", req.TopicId)
	// Pass empty string as title since it's not in the request
	id, err := h.postUC.CreatePost(ctx, req.TopicId, req.AuthorId, "", req.Content)
	if err != nil {
		h.logger.Error("failed to create post", "error", err)
		return nil, err
	}
	return &forumv1.CreatePostResponse{PostId: id}, nil
}

func (h *ForumHandler) GetPostByID(ctx context.Context, req *forumv1.GetPostRequest) (*forumv1.GetPostResponse, error) {
	post, err := h.postUC.GetPostByID(ctx, req.PostId)
	if err != nil {
		h.logger.Error("failed to get post", "error", err)
		return nil, err
	}
	return &forumv1.GetPostResponse{Post: toProtoPost(post)}, nil
}

func (h *ForumHandler) ListPostsByTopic(ctx context.Context, req *forumv1.ListPostsByTopicRequest) (*forumv1.ListPostsByTopicResponse, error) {
	limit := int(req.Limit)
	offset := int(req.Offset)
	posts, err := h.postUC.ListPostsByTopic(ctx, req.TopicId, limit, offset)
	if err != nil {
		h.logger.Error("failed to list posts", "error", err)
		return nil, err
	}
	protoPosts := make([]*forumv1.Post, len(posts))
	for i, p := range posts {
		protoPosts[i] = toProtoPost(p)
	}

	return &forumv1.ListPostsByTopicResponse{Posts: protoPosts}, nil
}

// Comment methods
func (h *ForumHandler) AddComment(ctx context.Context, req *forumv1.AddCommentRequest) (*forumv1.AddCommentResponse, error) {
	h.logger.Info("adding comment", "post_id", req.PostId)
	id, err := h.commentUC.CreateComment(ctx, req.PostId, req.AuthorId, req.Content)
	if err != nil {
		h.logger.Error("failed to create comment", "error", err)
		return nil, err
	}
	return &forumv1.AddCommentResponse{CommentId: id}, nil
}

func (h *ForumHandler) GetComment(ctx context.Context, req *forumv1.GetCommentRequest) (*forumv1.GetCommentResponse, error) {
	comment, err := h.commentUC.GetCommentByID(ctx, req.CommentId)
	if err != nil {
		h.logger.Error("failed to get comment", "error", err)
		return nil, err
	}
	return &forumv1.GetCommentResponse{Comment: toProtoComment(comment)}, nil
}

func (h *ForumHandler) ListCommentsByPost(ctx context.Context, req *forumv1.ListCommentsByPostRequest) (*forumv1.ListCommentsByPostResponse, error) {
	comments, err := h.commentUC.ListCommentsByPost(ctx, req.PostId)
	if err != nil {
		h.logger.Error("failed to list comments", "error", err)
		return nil, err
	}

	protoComments := make([]*forumv1.Comment, len(comments))
	for i, c := range comments {
		protoComments[i] = toProtoComment(c)
	}

	return &forumv1.ListCommentsByPostResponse{Comments: protoComments}, nil
}

// Tag methods
func (h *ForumHandler) CreateTag(ctx context.Context, req *forumv1.CreateTagRequest) (*forumv1.CreateTagResponse, error) {
	h.logger.Info("creating tag", "title", req.Title)
	id, err := h.tagUC.CreateTag(ctx, req.Title)
	if err != nil {
		h.logger.Error("failed to create tag", "error", err)
		return nil, err
	}
	return &forumv1.CreateTagResponse{TagId: id}, nil
}

func (h *ForumHandler) GetTagByID(ctx context.Context, req *forumv1.GetTagByIDRequest) (*forumv1.GetTagByIDResponse, error) {
	tag, err := h.tagUC.GetTagByID(ctx, req.TagId)
	if err != nil {
		h.logger.Error("failed to get tag", "error", err)
		return nil, err
	}
	return &forumv1.GetTagByIDResponse{Tag: toProtoTag(tag)}, nil
}

// Helper methods to convert domain entities to proto messages
func toProtoCategory(c *entity.Category) *forumv1.Category {
	return &forumv1.Category{
		Id:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func toProtoTopic(t *entity.Topic) *forumv1.Topic {
	return &forumv1.Topic{
		Id:         t.ID,
		Title:      t.Title,
		CategoryId: t.CategoryID,
		AuthorId:   t.AuthorID,
		CreatedAt:  t.CreatedAt,
	}
}

func toProtoPost(p *entity.Post) *forumv1.Post {
	return &forumv1.Post{
		Id:        p.ID,
		TopicId:   p.TopicID,
		AuthorId:  p.AuthorID,
		Content:   p.Content,
		Title:     p.Title,
		CreatedAt: p.CreatedAt,
	}
}

func toProtoComment(c *entity.Comment) *forumv1.Comment {
	return &forumv1.Comment{
		Id:        c.ID,
		PostId:    c.PostID,
		AuthorId:  c.AuthorID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
	}
}

func toProtoTag(t *entity.Tag) *forumv1.Tag {
	return &forumv1.Tag{
		Id:    t.ID,
		Title: t.Title,
	}
}
