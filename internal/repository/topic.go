package repository

import (
	"context"

	"github.com/VaneZ444/forum-service/internal/entity"
)

type TopicRepository interface {
	Create(ctx context.Context, topic *entity.Topic) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Topic, error)
	ListByCategory(ctx context.Context, categoryID int64, limit, offset int) ([]*entity.Topic, error)
	Delete(ctx context.Context, id int64) error
	UpdateTopic(ctx context.Context, topic *entity.Topic) error
}
