package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
	forumv1 "github.com/VaneZ444/golang-forum-protos/gen/go/forum"
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
		INSERT INTO topics (title, author_id, category_id, created_at, last_activity)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err = tx.QueryRowContext(ctx, topicQuery,
		topic.Title,
		topic.AuthorID,
		topic.CategoryID,
		topic.CreatedAt,
		topic.CreatedAt,
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

func (r *TopicRepository) List(ctx context.Context, categoryID *int64, limit, offset int, sorting *forumv1.Sorting) ([]*entity.Topic, int64, error) {
	query := `
		SELECT id, title, author_id, category_id, created_at, 
		       posts_count, views_count, last_activity, status
		FROM topics
	`
	countQuery := `SELECT COUNT(*) FROM topics`

	var args []any
	argIndex := 1

	if categoryID != nil {
		query += fmt.Sprintf(" WHERE category_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" WHERE category_id = $%d", argIndex)
		args = append(args, *categoryID)
		argIndex++
	}

	orderBy := buildOrderBy(sorting)
	query += fmt.Sprintf(" ORDER BY %s LIMIT $%d OFFSET $%d", orderBy, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list topics: %w", err)
	}
	defer rows.Close()

	topics := []*entity.Topic{}
	for rows.Next() {
		t := &entity.Topic{}
		if err := rows.Scan(
			&t.ID, &t.Title, &t.AuthorID, &t.CategoryID, &t.CreatedAt,
			&t.PostsCount, &t.ViewsCount, &t.LastActivity, &t.Status,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan topic: %w", err)
		}
		topics = append(topics, t)
	}

	var total int64
	err = r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get topic count: %w", err)
	}

	return topics, total, nil
}

func (r *TopicRepository) Update(ctx context.Context, topic *entity.Topic) (*entity.Topic, error) {
	const query = `
		UPDATE topics
		SET title = $1, category_id = $2, last_activity = $3
		WHERE id = $4
		RETURNING id, title, author_id, category_id, created_at, status, posts_count, views_count, last_activity
	`

	updatedTopic := &entity.Topic{}
	err := r.db.QueryRowContext(ctx, query,
		topic.Title,
		topic.CategoryID,
		topic.LastActivity,
		topic.ID,
	).Scan(
		&updatedTopic.ID,
		&updatedTopic.Title,
		&updatedTopic.AuthorID,
		&updatedTopic.CategoryID,
		&updatedTopic.CreatedAt,
		&updatedTopic.Status,
		&updatedTopic.PostsCount,
		&updatedTopic.ViewsCount,
		&updatedTopic.LastActivity,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update topic: %w", err)
	}

	return updatedTopic, nil
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

func buildOrderBy(sorting *forumv1.Sorting) string {
	if sorting == nil {
		return "last_activity DESC" // дефолт
	}

	var field string
	switch sorting.SortField {
	case forumv1.SortField_SORT_FIELD_CREATED_AT:
		field = "created_at"
	case forumv1.SortField_SORT_FIELD_UPDATED_AT:
		field = "last_activity"
	case forumv1.SortField_SORT_FIELD_TITLE:
		field = "title"
	case forumv1.SortField_SORT_FIELD_POPULARITY:
		field = "views_count" // можно по логике менять на posts_count
	default:
		field = "last_activity"
	}

	var order string
	switch sorting.SortOrder {
	case forumv1.SortOrder_SORT_ORDER_ASC:
		order = "ASC"
	case forumv1.SortOrder_SORT_ORDER_DESC:
		order = "DESC"
	default:
		order = "DESC"
	}

	return fmt.Sprintf("%s %s", field, order)
}
func (r *TopicRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Topic, int64, error) {
	tsquery := fmt.Sprintf("%s:*", strings.Join(strings.Fields(query), " & "))

	var total int64
	countQuery := `SELECT COUNT(*) FROM topics WHERE search_vector @@ to_tsquery('english', $1)`
	if err := r.db.QueryRowContext(ctx, countQuery, tsquery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count topics: %w", err)
	}

	searchQuery := `
		SELECT id, title, author_id, category_id, created_at, status, posts_count, views_count, last_activity
		FROM topics
		WHERE search_vector @@ to_tsquery('english', $1)
		ORDER BY ts_rank_cd(search_vector, to_tsquery('english', $1)) DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, searchQuery, tsquery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search topics: %w", err)
	}
	defer rows.Close()

	var topics []*entity.Topic
	for rows.Next() {
		var t entity.Topic
		if err := rows.Scan(
			&t.ID, &t.Title, &t.AuthorID, &t.CategoryID, &t.CreatedAt,
			&t.Status, &t.PostsCount, &t.ViewsCount, &t.LastActivity,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan topic: %w", err)
		}
		topics = append(topics, &t)
	}

	return topics, total, nil
}
