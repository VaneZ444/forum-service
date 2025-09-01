package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
	forumv1 "github.com/VaneZ444/golang-forum-protos/gen/go/forum"
)

type TopicUseCase interface {
	CreateTopic(ctx context.Context, topic *entity.Topic, post *entity.Post) (int64, int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Topic, *entity.Post, error)
	List(ctx context.Context, categoryID *int64, limit, offset int, sorting *forumv1.Sorting) ([]*entity.Topic, int64, error)
	UpdateTopic(ctx context.Context, topic *entity.Topic) (*entity.Topic, error)
	DeleteTopic(ctx context.Context, id int64) error
}

type topicUseCase struct {
	topicRepo    repository.TopicRepository
	categoryRepo repository.CategoryRepository
	logger       *slog.Logger
}

func NewTopicUseCase(
	topicRepo repository.TopicRepository,
	categoryRepo repository.CategoryRepository,
	logger *slog.Logger,
) TopicUseCase {
	return &topicUseCase{
		topicRepo:    topicRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *topicUseCase) CreateTopic(ctx context.Context, topic *entity.Topic, post *entity.Post) (int64, int64, error) {
	// Validate category exists
	_, err := uc.categoryRepo.GetByID(ctx, topic.CategoryID)
	if err != nil {
		uc.logger.Warn("category not found",
			slog.Int64("category_id", topic.CategoryID),
			slog.String("error", err.Error()),
		)
		return 0, 0, ErrCategoryNotFound
	}

	// Set timestamps
	now := time.Now().UTC()
	topic.CreatedAt = now
	post.CreatedAt = now

	// Create topic with first post
	err = uc.topicRepo.CreateWithPost(ctx, topic, post)
	if err != nil {
		uc.logger.Error("failed to create topic with post",
			slog.String("error", err.Error()),
		)
		return 0, 0, err
	}

	return topic.ID, post.ID, nil
}

func (uc *topicUseCase) GetByID(ctx context.Context, id int64) (*entity.Topic, *entity.Post, error) {
	topic, post, err := uc.topicRepo.GetByIDWithFirstPost(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			uc.logger.Warn("topic not found", slog.Int64("id", id))
			return nil, nil, ErrTopicNotFound
		}
		uc.logger.Error("failed to get topic",
			slog.Int64("id", id),
			slog.String("error", err.Error()),
		)
		return nil, nil, err
	}
	return topic, post, nil
}

func (uc *topicUseCase) List(
	ctx context.Context,
	categoryID *int64,
	limit, offset int,
	sorting *forumv1.Sorting,
) ([]*entity.Topic, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	// Проверяем категорию, если указана
	if categoryID != nil {
		_, err := uc.categoryRepo.GetByID(ctx, *categoryID)
		if err != nil {
			uc.logger.Warn("category not found",
				slog.Int64("category_id", *categoryID),
				slog.String("error", err.Error()),
			)
			return nil, 0, ErrCategoryNotFound
		}
	}

	// Передаем сортировку дальше в репозиторий
	return uc.topicRepo.List(ctx, categoryID, limit, offset, sorting)
}

func (uc *topicUseCase) UpdateTopic(ctx context.Context, topic *entity.Topic) (*entity.Topic, error) {
	existing, err := uc.topicRepo.GetByID(ctx, topic.ID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	// Preserve неизменяемые поля
	topic.AuthorID = existing.AuthorID
	topic.CreatedAt = existing.CreatedAt
	topic.Status = existing.Status
	topic.PostsCount = existing.PostsCount
	topic.ViewsCount = existing.ViewsCount

	// Обновляем last_activity
	topic.LastActivity = time.Now().UTC()

	return uc.topicRepo.Update(ctx, topic)
}

func (uc *topicUseCase) DeleteTopic(ctx context.Context, id int64) error {
	// Check if topic exists
	_, _, err := uc.topicRepo.GetByIDWithFirstPost(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return ErrTopicNotFound
		}
		return err
	}

	return uc.topicRepo.Delete(ctx, id)
}
