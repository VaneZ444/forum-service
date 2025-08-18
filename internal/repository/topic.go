package repository

import (
	"context"

	"github.com/VaneZ444/forum-service/internal/entity"
)

type TopicRepository interface {
	CreateWithPost(ctx context.Context, topic *entity.Topic, post *entity.Post) error
	GetByID(ctx context.Context, id int64) (*entity.Topic, error)
	GetByIDWithFirstPost(ctx context.Context, id int64) (*entity.Topic, *entity.Post, error)
	List(ctx context.Context, categoryID int64, limit, offset int) ([]*entity.Topic, int64, error)
	Update(ctx context.Context, topic *entity.Topic) error
	Delete(ctx context.Context, id int64) error
}
