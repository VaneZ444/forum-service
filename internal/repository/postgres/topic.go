package postgres

import (
	"context"
	"database/sql"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

type TopicRepository struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) repository.TopicRepository {
	return &TopicRepository{db: db}
}

func (r *TopicRepository) Create(ctx context.Context, topic *entity.Topic) (int64, error) {
	query := `INSERT INTO topics (title, author_id, category_id, created_at)
	          VALUES ($1, $2, $3, $4)
	          RETURNING id`

	var id int64
	err := r.db.QueryRowContext(ctx, query,
		topic.Title,
		topic.AuthorID,
		topic.CategoryID,
		topic.CreatedAt,
	).Scan(&id)

	return id, err
}

func (r *TopicRepository) GetByID(ctx context.Context, id int64) (*entity.Topic, error) {
	query := `SELECT id, title, author_id, category_id, created_at FROM topics WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var topic entity.Topic
	err := row.Scan(
		&topic.ID,
		&topic.Title,
		&topic.AuthorID,
		&topic.CategoryID,
		&topic.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &topic, nil
}

func (r *TopicRepository) ListByCategory(ctx context.Context, categoryID int64, limit, offset int) ([]*entity.Topic, error) {
	query := `SELECT id, title, author_id, category_id, created_at
	          FROM topics
	          WHERE category_id = $1
	          ORDER BY created_at DESC
	          LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []*entity.Topic
	for rows.Next() {
		var topic entity.Topic
		err := rows.Scan(
			&topic.ID,
			&topic.Title,
			&topic.AuthorID,
			&topic.CategoryID,
			&topic.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		topics = append(topics, &topic)
	}

	return topics, nil
}

func (r *TopicRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM topics WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
