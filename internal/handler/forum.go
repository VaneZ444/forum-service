package handler

import (
	"context"
	"errors"
	"log/slog"

	forumv1 "github.com/VaneZ444/golang-forum-protos/gen/go/forum"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/usecase"
)

type ForumHandler struct {
	forumv1.UnimplementedForumServer

	topicUC usecase.TopicUseCase
	logger  *slog.Logger
}

func New(topicUC usecase.TopicUseCase, logger *slog.Logger) *ForumHandler {
	return &ForumHandler{
		topicUC: topicUC,
		logger:  logger,
	}
}

func (h *ForumHandler) CreateTopic(ctx context.Context, req *forumv1.CreateTopicRequest) (*forumv1.CreateTopicResponse, error) {
	id, err := h.topicUC.CreateTopic(ctx, req.Title, req.AuthorId, req.CategoryId)
	if err != nil {
		h.logger.Warn("failed to create topic", slog.String("err", err.Error()))
		return nil, err
	}

	return &forumv1.CreateTopicResponse{
		TopicId: id,
	}, nil
}

func (h *ForumHandler) GetTopicByID(ctx context.Context, req *forumv1.GetTopicByIDRequest) (*forumv1.GetTopicByIDResponse, error) {
	topic, err := h.topicUC.GetTopicByID(ctx, req.TopicId)
	if err != nil {
		if errors.Is(err, usecase.ErrTopicNotFound) {
			return nil, err // здесь можно добавить status.Error если нужно
		}
		h.logger.Error("failed to get topic", slog.String("err", err.Error()))
		return nil, err
	}

	return &forumv1.GetTopicByIDResponse{
		Topic: toProtoTopic(topic),
	}, nil
}

func (h *ForumHandler) ListTopicsByCategory(ctx context.Context, req *forumv1.ListTopicsByCategoryRequest) (*forumv1.ListTopicsByCategoryResponse, error) {
	topics, err := h.topicUC.ListTopicsByCategory(ctx, req.CategoryId, int(req.Limit), int(req.Offset))

	if err != nil {
		h.logger.Error("failed to list topics", slog.String("err", err.Error()))
		return nil, err
	}

	var protoTopics []*forumv1.Topic
	for _, topic := range topics {
		protoTopics = append(protoTopics, toProtoTopic(topic))
	}

	return &forumv1.ListTopicsByCategoryResponse{
		Topics: protoTopics,
	}, nil
}
func toProtoTopic(t *entity.Topic) *forumv1.Topic {
	return &forumv1.Topic{
		Id:         t.ID,
		Title:      t.Title,
		AuthorId:   t.AuthorID,
		CategoryId: t.CategoryID,
		CreatedAt:  t.CreatedAt,
	}
}
