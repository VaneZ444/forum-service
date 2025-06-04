package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

type TopicUseCase interface {
	CreateTopic(ctx context.Context, title string, authorID int64, categoryID int64) (int64, error)
	GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error)
	ListTopicsByCategory(ctx context.Context, categoryID int64, limit, offset int) ([]*entity.Topic, error)
}

type topicUseCase struct {
	topicRepo    repository.TopicRepository
	categoryRepo repository.CategoryRepository
	logger       *slog.Logger
}

func NewTopicUseCase(topicRepo repository.TopicRepository, categoryRepo repository.CategoryRepository, logger *slog.Logger) TopicUseCase {
	return &topicUseCase{
		topicRepo:    topicRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *topicUseCase) CreateTopic(ctx context.Context, title string, authorID int64, categoryID int64) (int64, error) {
	_, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		uc.logger.Warn("failed to find category", slog.Int64("categoryID", categoryID), slog.String("err", err.Error()))
		return 0, ErrCategoryNotFound
	}

	topic := &entity.Topic{
		Title:      title,
		AuthorID:   authorID,
		CategoryID: categoryID,
		CreatedAt:  time.Now().Unix(),
	}

	id, err := uc.topicRepo.Create(ctx, topic)
	if err != nil {
		uc.logger.Error("failed to create topic", slog.String("err", err.Error()))
		return 0, err
	}

	return id, nil
}

func (uc *topicUseCase) GetTopicByID(ctx context.Context, id int64) (*entity.Topic, error) {
	topic, err := uc.topicRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Warn("topic not found", slog.Int64("id", id), slog.String("err", err.Error()))
		return nil, ErrTopicNotFound
	}
	return topic, nil
}

func (uc *topicUseCase) ListTopicsByCategory(ctx context.Context, categoryID int64, limit, offset int) ([]*entity.Topic, error) {
	return uc.topicRepo.ListByCategory(ctx, categoryID, limit, offset)
}
