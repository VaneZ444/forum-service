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
