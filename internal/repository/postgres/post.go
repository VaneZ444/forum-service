package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) repository.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *entity.Post) (int64, error) {
	const query = `INSERT INTO posts (topic_id, content, author_id, created_at) 
				   VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, post.TopicID, post.Content, post.AuthorID, post.CreatedAt).Scan(&post.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to create post: %w", err)
	}

	return post.ID, nil
}

func (r *postRepository) GetByID(ctx context.Context, id int64) (*entity.Post, error) {
	const query = `SELECT id, topic_id, content, author_id, created_at FROM posts WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var post entity.Post
	if err := row.Scan(&post.ID, &post.TopicID, &post.Content, &post.AuthorID, &post.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}

	return &post, nil
}

func (r *postRepository) Update(ctx context.Context, post *entity.Post) error {
	query := `UPDATE posts 
              SET title = $1, content = $2
              WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found with id: %d", post.ID)
	}

	return nil
}

func (r *postRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found with id: %d", id)
	}

	return nil
}

func (r *postRepository) ListByTag(ctx context.Context, tagID int64, limit, offset int) ([]*entity.Post, error) {
	query := `SELECT p.id, p.topic_id, p.title, p.content, p.author_id, p.created_at 
              FROM posts p
              JOIN post_tags pt ON p.id = pt.post_id
              WHERE pt.tag_id = $1
              ORDER BY p.created_at DESC
              LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, tagID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts by tag: %w", err)
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		var p entity.Post
		if err := rows.Scan(&p.ID, &p.TopicID, &p.Title, &p.Content, &p.AuthorID, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, &p)
	}
	return posts, nil
}

func (r *postRepository) ListByTopic(ctx context.Context, topicID int64, limit int, offset int) ([]*entity.Post, error) {
	const query = `SELECT id, topic_id, content, author_id, created_at 
				   FROM posts WHERE topic_id = $1 
				   ORDER BY created_at ASC
				   LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, topicID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		var p entity.Post
		if err := rows.Scan(&p.ID, &p.TopicID, &p.Content, &p.AuthorID, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return posts, nil
}

func (r *postRepository) List(ctx context.Context, topicID, tagID int64, limit, offset int) ([]*entity.Post, int64, error) {
	query := `SELECT id, title, content, author_id, topic_id, created_at, updated_at 
              FROM posts WHERE 1=1`
	args := []interface{}{}
	idx := 1

	if topicID > 0 {
		query += fmt.Sprintf(" AND topic_id = $%d", idx)
		args = append(args, topicID)
		idx++
	}
	if tagID > 0 {
		query += fmt.Sprintf(" AND id IN (SELECT post_id FROM post_tags WHERE tag_id = $%d)", idx)
		args = append(args, tagID)
		idx++
	}

	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS count_subquery"
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", idx, idx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list posts: %w", err)
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		p := new(entity.Post)
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.TopicID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}

	return posts, total, nil
}
