package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

type TopicRepository struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) repository.TopicRepository {
	return &TopicRepository{db: db}
}

func (r *TopicRepository) CreateWithPost(ctx context.Context, topic *entity.Topic, post *entity.Post) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert topic
	topicQuery := `
		INSERT INTO topics (title, author_id, category_id, created_at, posts_count, last_activity)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err = tx.QueryRowContext(ctx, topicQuery,
		topic.Title,
		topic.AuthorID,
		topic.CategoryID,
		topic.CreatedAt,
		1,               // Initial posts_count
		topic.CreatedAt, // Last activity same as creation time
	).Scan(&topic.ID)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	// Insert first post
	post.TopicID = topic.ID
	postQuery := `
		INSERT INTO posts (topic_id, author_id, title, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err = tx.QueryRowContext(ctx, postQuery,
		post.TopicID,
		post.AuthorID,
		post.Title,
		post.Content,
		post.CreatedAt,
	).Scan(&post.ID)
	if err != nil {
		return fmt.Errorf("failed to create first post: %w", err)
	}

	// Update category topic count
	updateCategoryQuery := `
		UPDATE categories 
		SET topics_count = topics_count + 1 
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, updateCategoryQuery, topic.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to update category topic count: %w", err)
	}

	return tx.Commit()
}
func (r *TopicRepository) GetByID(ctx context.Context, id int64) (*entity.Topic, error) {
	query := `
		SELECT 
			id, title, author_id, category_id, created_at, 
			posts_count, views_count, last_activity, status
		FROM topics
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	topic := &entity.Topic{}
	err := row.Scan(
		&topic.ID, &topic.Title, &topic.AuthorID, &topic.CategoryID, &topic.CreatedAt,
		&topic.PostsCount, &topic.ViewsCount, &topic.LastActivity, &topic.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get topic by id: %w", err)
	}

	return topic, nil
}

func (r *TopicRepository) GetByIDWithFirstPost(ctx context.Context, id int64) (*entity.Topic, *entity.Post, error) {
	query := `
		SELECT 
			t.id, t.title, t.author_id, t.category_id, t.created_at, 
			t.posts_count, t.views_count, t.last_activity, t.status,
			p.id, p.author_id, p.title, p.content, p.created_at
		FROM topics t
		JOIN posts p ON t.id = p.topic_id
		WHERE t.id = $1
		ORDER BY p.created_at ASC
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	topic := &entity.Topic{}
	post := &entity.Post{}

	err := row.Scan(
		&topic.ID, &topic.Title, &topic.AuthorID, &topic.CategoryID, &topic.CreatedAt,
		&topic.PostsCount, &topic.ViewsCount, &topic.LastActivity, &topic.Status,
		&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, repository.ErrNotFound
		}
		return nil, nil, fmt.Errorf("failed to get topic with first post: %w", err)
	}

	post.TopicID = topic.ID
	return topic, post, nil
}

func (r *TopicRepository) List(ctx context.Context, categoryID int64, limit, offset int) ([]*entity.Topic, int64, error) {
	query := `
		SELECT id, title, author_id, category_id, created_at, 
		       posts_count, views_count, last_activity, status
		FROM topics
		WHERE category_id = $1
		ORDER BY last_activity DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, categoryID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list topics: %w", err)
	}
	defer rows.Close()

	topics := []*entity.Topic{}
	for rows.Next() {
		t := &entity.Topic{}
		err := rows.Scan(
			&t.ID, &t.Title, &t.AuthorID, &t.CategoryID, &t.CreatedAt,
			&t.PostsCount, &t.ViewsCount, &t.LastActivity, &t.Status,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan topic: %w", err)
		}
		topics = append(topics, t)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM topics WHERE category_id = $1`
	var total int64
	err = r.db.QueryRowContext(ctx, countQuery, categoryID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get topic count: %w", err)
	}

	return topics, total, nil
}

func (r *TopicRepository) Update(ctx context.Context, topic *entity.Topic) error {
	query := `
		UPDATE topics 
		SET title = $1, category_id = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query,
		topic.Title,
		topic.CategoryID,
		time.Now().UnixMilli(),
		topic.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update topic: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *TopicRepository) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get topic to update category count
	topicQuery := `SELECT category_id FROM topics WHERE id = $1`
	var categoryID int64
	err = tx.QueryRowContext(ctx, topicQuery, id).Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.ErrNotFound
		}
		return fmt.Errorf("failed to get topic: %w", err)
	}

	// Delete topic
	deleteQuery := `DELETE FROM topics WHERE id = $1`
	result, err := tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	// Update category topic count
	updateCategoryQuery := `
		UPDATE categories 
		SET topics_count = topics_count - 1 
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, updateCategoryQuery, categoryID)
	if err != nil {
		return fmt.Errorf("failed to update category topic count: %w", err)
	}

	return tx.Commit()
}
