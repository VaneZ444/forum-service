package handler

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/usecase"
	forumv1 "github.com/VaneZ444/golang-forum-protos/gen/go/forum"
	"github.com/gosimple/slug"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	createdCategory, err := h.categoryUC.CreateCategory(ctx, category)
	if err != nil {
		h.logger.Error("failed to create category", "error", err)
		return nil, err
	}
	return &forumv1.CategoryResponse{Category: toProtoCategory(createdCategory)}, nil
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
	h.logger.Info("updating category", "id", req.GetId())

	// 1) Берём текущую версию
	existing, err := h.categoryUC.GetByID(ctx, req.GetId())
	if err != nil {
		h.logger.Error("get category failed", "error", err)
		return nil, err
	}

	// 2) Мержим изменения из запроса
	title := existing.Title
	if req.Title != nil {
		t := strings.TrimSpace(req.GetTitle())
		if t == "" {
			return nil, status.Error(codes.InvalidArgument, "title cannot be empty")
		}
		title = t
	}

	description := existing.Description
	if req.Description != nil {
		description = req.GetDescription()
	}

	// 3) Генерируем новый slug, если title поменялся
	newSlug := existing.Slug
	if title != existing.Title {
		newSlug = slug.Make(title)

		// Проверяем уникальность
		if other, _ := h.categoryUC.GetBySlug(ctx, newSlug); other != nil && other.ID != existing.ID {
			return nil, status.Errorf(codes.AlreadyExists, "slug %s already exists", newSlug)
		}
	}

	// 4) Собираем сущность для апдейта
	cat := &entity.Category{
		ID:          existing.ID,
		Title:       title,
		Slug:        newSlug,
		Description: description,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   time.Now().UTC(),
	}

	// 5) Апдейт
	updated, err := h.categoryUC.UpdateCategory(ctx, cat)
	if err != nil {
		h.logger.Error("update category failed", "error", err)
		return nil, err
	}

	return &forumv1.CategoryResponse{Category: toProtoCategory(updated)}, nil
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

func (h *ForumHandler) UpdateTopic(ctx context.Context, req *forumv1.UpdateTopicRequest) (*forumv1.TopicResponse, error) {
	h.logger.Info("updating topic", "id", req.GetId())

	// 1) Берём текущий топик
	existing, _, err := h.topicUC.GetByID(ctx, req.GetId())
	if err != nil {
		h.logger.Error("get topic failed", "error", err)
		return nil, err
	}

	// 2) Мержим изменения
	title := existing.Title
	if req.Title != nil {
		t := strings.TrimSpace(req.GetTitle())
		if t == "" {
			return nil, status.Error(codes.InvalidArgument, "title cannot be empty")
		}
		title = t
	}

	categoryID := existing.CategoryID
	if req.CategoryId != nil {
		categoryID = req.GetCategoryId()
	}

	// 3) Собираем обновлённую сущность
	topic := &entity.Topic{
		ID:           existing.ID,
		Title:        title,
		CategoryID:   categoryID,
		AuthorID:     existing.AuthorID,
		CreatedAt:    existing.CreatedAt,
		Status:       existing.Status,
		PostsCount:   existing.PostsCount,
		ViewsCount:   existing.ViewsCount,
		LastActivity: time.Now().UTC(), // обновляем активность
	}

	// 4) Апдейт
	updated, err := h.topicUC.UpdateTopic(ctx, topic)
	if err != nil {
		h.logger.Error("update topic failed", "error", err)
		return nil, err
	}

	return &forumv1.TopicResponse{
		Topic:     toProtoTopic(updated),
		FirstPost: nil, // посты не трогаем
	}, nil
}

func (h *ForumHandler) ListTopics(ctx context.Context, req *forumv1.ListTopicsRequest) (*forumv1.ListTopicsResponse, error) {
	var categoryID *int64
	if req.CategoryId != nil {
		id := req.GetCategoryId()
		categoryID = &id
	}

	// Пагинация
	limit := 50
	offset := 0
	if req.Pagination != nil {
		if req.Pagination.Limit > 0 {
			limit = int(req.Pagination.Limit)
		}
		if req.Pagination.Offset > 0 {
			offset = int(req.Pagination.Offset)
		}
	}

	// Вызываем юзкейс
	topics, total, err := h.topicUC.List(ctx, categoryID, limit, offset, req.GetSorting())
	if err != nil {
		h.logger.Error("failed to list topics", "error", err)
		return nil, err
	}

	// Конвертируем в proto
	protoTopics := make([]*forumv1.Topic, len(topics))
	for i, t := range topics {
		protoTopics[i] = toProtoTopic(t)
	}

	return &forumv1.ListTopicsResponse{
		Topics:     protoTopics,
		TotalCount: total,
	}, nil
}

func (h *ForumHandler) DeleteTopic(ctx context.Context, req *forumv1.DeleteTopicRequest) (*forumv1.Empty, error) {
	h.logger.Info("deleting topic", "id", req.GetId())

	err := h.topicUC.DeleteTopic(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, usecase.ErrTopicNotFound) {
			return nil, status.Error(codes.NotFound, "topic not found")
		}
		h.logger.Error("failed to delete topic", "error", err)
		return nil, status.Error(codes.Internal, "failed to delete topic")
	}

	return &forumv1.Empty{}, nil
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
	userID := GetUserIDFromCtx(ctx)
	if err := h.postUC.AddView(ctx, req.GetId(), userID); err != nil {
		h.logger.Warn("failed to add post view", "error", err)
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

func (h *ForumHandler) UpdatePost(ctx context.Context, req *forumv1.UpdatePostRequest) (*forumv1.PostResponse, error) {
	h.logger.Info("updating post", "id", req.GetId())

	post, err := h.postUC.UpdatePost(ctx, req)
	if err != nil {
		if errors.Is(err, usecase.ErrPostNotFound) {
			return nil, status.Error(codes.NotFound, "post not found")
		}
		return nil, status.Error(codes.Internal, "failed to update post")
	}

	return &forumv1.PostResponse{Post: toProtoPost(post)}, nil
}

func (h *ForumHandler) DeletePost(ctx context.Context, req *forumv1.DeletePostRequest) (*forumv1.Empty, error) {
	h.logger.Info("deleting post", "id", req.GetId())

	err := h.postUC.DeletePost(ctx, req.GetId())
	if err != nil {
		h.logger.Error("failed to delete post", "error", err)
		return nil, err
	}

	return &forumv1.Empty{}, nil
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
		CreatedAt:   timestamppb.New(c.CreatedAt),
		UpdatedAt:   timestamppb.New(c.UpdatedAt),
	}
}

func toProtoTopic(t *entity.Topic) *forumv1.Topic {
	return &forumv1.Topic{
		Id:           t.ID,
		Title:        t.Title,
		AuthorId:     t.AuthorID,
		CategoryId:   t.CategoryID,
		CreatedAt:    timestamppb.New(t.CreatedAt),
		Status:       forumv1.Status(t.Status),
		PostsCount:   t.PostsCount,
		ViewsCount:   t.ViewsCount,
		LastActivity: timestamppb.New(t.LastActivity),
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
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timestamppb.New(p.UpdatedAt),
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
		CreatedAt: timestamppb.New(c.CreatedAt),
		UpdatedAt: timestamppb.New(c.UpdatedAt),
	}
}

func toProtoTag(t *entity.Tag) *forumv1.Tag {
	return &forumv1.Tag{
		Id:   t.ID,
		Name: t.Name,
		Slug: t.Slug,
	}
}
