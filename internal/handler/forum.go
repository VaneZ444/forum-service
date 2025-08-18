package handler

import (
	"context"
	"log/slog"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/usecase"
	forumv1 "github.com/VaneZ444/golang-forum-protos/gen/go/forum"
)

type ForumHandler struct {
	forumv1.UnimplementedForumServiceServer
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

// ================== Category Handlers ==================
func (h *ForumHandler) CreateCategory(ctx context.Context, req *forumv1.CreateCategoryRequest) (*forumv1.CategoryResponse, error) {
	h.logger.Info("creating category", "title", req.GetTitle())
	category := &entity.Category{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
	}

	err := h.categoryUC.CreateCategory(ctx, category)
	if err != nil {
		h.logger.Error("failed to create category", "error", err)
		return nil, err
	}
	return &forumv1.CategoryResponse{Category: toProtoCategory(category)}, nil
}

func (h *ForumHandler) GetCategory(ctx context.Context, req *forumv1.GetCategoryRequest) (*forumv1.CategoryResponse, error) {
	category, err := h.categoryUC.GetByID(ctx, req.GetId())
	if err != nil {
		h.logger.Error("failed to get category", "error", err)
		return nil, err
	}
	return &forumv1.CategoryResponse{Category: toProtoCategory(category)}, nil
}

func (h *ForumHandler) ListCategories(ctx context.Context, req *forumv1.ListCategoriesRequest) (*forumv1.ListCategoriesResponse, error) {
	pagination := req.GetPagination()
	limit := 100
	offset := 0
	if pagination != nil {
		limit = int(pagination.GetLimit())
		offset = int(pagination.GetOffset())
	}

	categories, total, err := h.categoryUC.List(ctx, limit, offset)
	if err != nil {
		h.logger.Error("failed to list categories", "error", err)
		return nil, err
	}

	protoCategories := make([]*forumv1.Category, len(categories))
	for i, c := range categories {
		protoCategories[i] = toProtoCategory(c)
	}

	return &forumv1.ListCategoriesResponse{
		Categories: protoCategories,
		TotalCount: total,
	}, nil
}

func (h *ForumHandler) UpdateCategory(ctx context.Context, req *forumv1.UpdateCategoryRequest) (*forumv1.CategoryResponse, error) {
	// Implementation will be similar to CreateCategory
	return nil, nil
}

func (h *ForumHandler) DeleteCategory(ctx context.Context, req *forumv1.DeleteCategoryRequest) (*forumv1.Empty, error) {
	err := h.categoryUC.DeleteCategory(ctx, req.GetId())
	if err != nil {
		h.logger.Error("failed to delete category", "error", err)
		return nil, err
	}
	return &forumv1.Empty{}, nil
}

// ================== Topic Handlers ==================
func (h *ForumHandler) CreateTopic(ctx context.Context, req *forumv1.CreateTopicRequest) (*forumv1.TopicResponse, error) {
	h.logger.Info("creating topic", "title", req.GetTitle())

	// Create topic entity
	topic := &entity.Topic{
		Title:      req.GetTitle(),
		AuthorID:   req.GetAuthorId(),
		CategoryID: req.GetCategoryId(),
	}

	// Create first post entity
	post := &entity.Post{
		Title:    req.GetTitle(),
		Content:  req.GetContent(),
		AuthorID: req.GetAuthorId(),
	}

	topicID, postID, err := h.topicUC.CreateTopic(ctx, topic, post)
	if err != nil {
		h.logger.Error("failed to create topic", "error", err)
		return nil, err
	}

	// Set IDs for response
	topic.ID = topicID
	post.ID = postID
	post.TopicID = topicID

	return &forumv1.TopicResponse{
		Topic:     toProtoTopic(topic),
		FirstPost: toProtoPost(post),
	}, nil
}

func (h *ForumHandler) GetTopic(ctx context.Context, req *forumv1.GetTopicRequest) (*forumv1.TopicResponse, error) {
	topic, firstPost, err := h.topicUC.GetByID(ctx, req.GetId())
	if err != nil {
		h.logger.Error("failed to get topic", "error", err)
		return nil, err
	}
	return &forumv1.TopicResponse{
		Topic:     toProtoTopic(topic),
		FirstPost: toProtoPost(firstPost),
	}, nil
}

func (h *ForumHandler) ListTopics(ctx context.Context, req *forumv1.ListTopicsRequest) (*forumv1.ListTopicsResponse, error) {
	var categoryID int64
	if req.CategoryId != nil {
		categoryID = req.GetCategoryId()
	}

	pagination := req.GetPagination()
	limit := 50
	offset := 0
	if pagination != nil {
		limit = int(pagination.GetLimit())
		offset = int(pagination.GetOffset())
	}

	topics, total, err := h.topicUC.List(ctx, categoryID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list topics", "error", err)
		return nil, err
	}

	protoTopics := make([]*forumv1.Topic, len(topics))
	for i, t := range topics {
		protoTopics[i] = toProtoTopic(t)
	}

	return &forumv1.ListTopicsResponse{
		Topics:     protoTopics,
		TotalCount: total,
	}, nil
}

// ================== Post Handlers ==================
func (h *ForumHandler) CreatePost(ctx context.Context, req *forumv1.CreatePostRequest) (*forumv1.PostResponse, error) {
	h.logger.Info("creating post", "topic_id", req.GetTopicId())

	post := &entity.Post{
		TopicID:  req.GetTopicId(),
		AuthorID: req.GetAuthorId(),
		Title:    req.GetTitle(),
		Content:  req.GetContent(),
		Images:   req.GetImages(),
	}

	// Add tags if provided
	if req.TagIds != nil {
		post.Tags = make([]entity.Tag, len(req.GetTagIds()))
		for i, tagID := range req.GetTagIds() {
			post.Tags[i] = entity.Tag{ID: tagID}
		}
	}

	id, err := h.postUC.CreatePost(ctx, post)
	if err != nil {
		h.logger.Error("failed to create post", "error", err)
		return nil, err
	}
	post.ID = id
	return &forumv1.PostResponse{Post: toProtoPost(post)}, nil
}

func (h *ForumHandler) GetPost(ctx context.Context, req *forumv1.GetPostRequest) (*forumv1.PostResponse, error) {
	post, err := h.postUC.GetPostByID(ctx, req.GetId())
	if err != nil {
		h.logger.Error("failed to get post", "error", err)
		return nil, err
	}
	return &forumv1.PostResponse{Post: toProtoPost(post)}, nil
}

func (h *ForumHandler) ListPosts(ctx context.Context, req *forumv1.ListPostsRequest) (*forumv1.ListPostsResponse, error) {
	var topicID, tagID int64
	if req.TopicId != nil {
		topicID = req.GetTopicId()
	}
	if req.TagId != nil {
		tagID = req.GetTagId()
	}

	pagination := req.GetPagination()
	limit := 50
	offset := 0
	if pagination != nil {
		limit = int(pagination.GetLimit())
		offset = int(pagination.GetOffset())
	}

	// Sorting handling would go here

	posts, total, err := h.postUC.List(ctx, topicID, tagID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list posts", "error", err)
		return nil, err
	}

	protoPosts := make([]*forumv1.Post, len(posts))
	for i, p := range posts {
		protoPosts[i] = toProtoPost(p)
	}

	return &forumv1.ListPostsResponse{
		Posts:      protoPosts,
		TotalCount: total,
	}, nil
}

// ================== Comment Handlers ==================
func (h *ForumHandler) CreateComment(ctx context.Context, req *forumv1.CreateCommentRequest) (*forumv1.CommentResponse, error) {
	h.logger.Info("creating comment", "post_id", req.GetPostId())

	comment := &entity.Comment{
		PostID:   req.GetPostId(),
		AuthorID: req.GetAuthorId(),
		Content:  req.GetContent(),
	}

	id, err := h.commentUC.CreateComment(ctx, comment)
	if err != nil {
		h.logger.Error("failed to create comment", "error", err)
		return nil, err
	}
	comment.ID = id

	return &forumv1.CommentResponse{Comment: toProtoComment(comment)}, nil
}

func (h *ForumHandler) GetComment(ctx context.Context, req *forumv1.GetCommentRequest) (*forumv1.CommentResponse, error) {
	comment, err := h.commentUC.GetCommentByID(ctx, req.GetId())
	if err != nil {
		h.logger.Error("failed to get comment", "error", err)
		return nil, err
	}
	return &forumv1.CommentResponse{Comment: toProtoComment(comment)}, nil
}

func (h *ForumHandler) ListComments(ctx context.Context, req *forumv1.ListCommentsRequest) (*forumv1.ListCommentsResponse, error) {
	pagination := req.GetPagination()
	limit := 100
	offset := 0
	if pagination != nil {
		limit = int(pagination.GetLimit())
		offset = int(pagination.GetOffset())
	}

	comments, total, err := h.commentUC.ListByPost(ctx, req.GetPostId(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list comments", "error", err)
		return nil, err
	}

	protoComments := make([]*forumv1.Comment, len(comments))
	for i, c := range comments {
		protoComments[i] = toProtoComment(c)
	}

	return &forumv1.ListCommentsResponse{
		Comments:   protoComments,
		TotalCount: total,
	}, nil
}

// ================== Tag Handlers ==================
func (h *ForumHandler) CreateTag(ctx context.Context, req *forumv1.CreateTagRequest) (*forumv1.TagResponse, error) {
	h.logger.Info("creating tag", "name", req.GetName())

	tag := &entity.Tag{
		Name: req.GetName(),
	}

	err := h.tagUC.CreateTag(ctx, tag)
	if err != nil {
		h.logger.Error("failed to create tag", "error", err)
		return nil, err
	}
	return &forumv1.TagResponse{Tag: toProtoTag(tag)}, nil
}

func (h *ForumHandler) GetTag(ctx context.Context, req *forumv1.GetTagRequest) (*forumv1.TagResponse, error) {
	var tag *entity.Tag
	var err error

	switch id := req.Identifier.(type) {
	case *forumv1.GetTagRequest_Id:
		tag, err = h.tagUC.GetTagByID(ctx, id.Id)
	case *forumv1.GetTagRequest_Slug:
		tag, err = h.tagUC.GetTagBySlug(ctx, id.Slug)
	default:
		h.logger.Error("invalid tag identifier")
		return nil, err
	}

	if err != nil {
		h.logger.Error("failed to get tag", "error", err)
		return nil, err
	}
	return &forumv1.TagResponse{Tag: toProtoTag(tag)}, nil
}

func (h *ForumHandler) ListTags(ctx context.Context, req *forumv1.ListTagsRequest) (*forumv1.ListTagsResponse, error) {
	pagination := req.GetPagination()
	limit := 100
	offset := 0
	if pagination != nil {
		limit = int(pagination.GetLimit())
		offset = int(pagination.GetOffset())
	}

	tags, total, err := h.tagUC.List(ctx, limit, offset)
	if err != nil {
		h.logger.Error("failed to list tags", "error", err)
		return nil, err
	}

	protoTags := make([]*forumv1.Tag, len(tags))
	for i, t := range tags {
		protoTags[i] = toProtoTag(t)
	}

	return &forumv1.ListTagsResponse{
		Tags:       protoTags,
		TotalCount: total,
	}, nil
}

// ================== Helper Functions ==================
func toProtoCategory(c *entity.Category) *forumv1.Category {
	return &forumv1.Category{
		Id:          c.ID,
		Title:       c.Title,
		Slug:        c.Slug,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func toProtoTopic(t *entity.Topic) *forumv1.Topic {
	return &forumv1.Topic{
		Id:           t.ID,
		Title:        t.Title,
		AuthorId:     t.AuthorID,
		CategoryId:   t.CategoryID,
		CreatedAt:    t.CreatedAt,
		Status:       forumv1.Status(t.Status),
		PostsCount:   t.PostsCount,
		ViewsCount:   t.ViewsCount,
		LastActivity: t.LastActivity,
	}
}

func toProtoPost(p *entity.Post) *forumv1.Post {
	tags := make([]*forumv1.Tag, len(p.Tags))
	for i, t := range p.Tags {
		tags[i] = toProtoTag(&t)
	}

	return &forumv1.Post{
		Id:            p.ID,
		TopicId:       p.TopicID,
		AuthorId:      p.AuthorID,
		Title:         p.Title,
		Content:       p.Content,
		Images:        p.Images,
		Tags:          tags,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
		Status:        forumv1.Status(p.Status),
		ViewsCount:    p.ViewsCount,
		CommentsCount: p.CommentsCount,
		LikesCount:    p.LikesCount,
	}
}

func toProtoComment(c *entity.Comment) *forumv1.Comment {
	return &forumv1.Comment{
		Id:        c.ID,
		PostId:    c.PostID,
		AuthorId:  c.AuthorID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toProtoTag(t *entity.Tag) *forumv1.Tag {
	return &forumv1.Tag{
		Id:   t.ID,
		Name: t.Name,
		Slug: t.Slug,
	}
}
