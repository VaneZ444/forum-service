package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) repository.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *entity.Comment) (int64, error) {
	const query = `INSERT INTO comments (post_id, content, author_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, comment.PostID, comment.Content, comment.AuthorID, comment.CreatedAt).Scan(&comment.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to create comment: %w", err)
	}
	return comment.ID, nil
}

func (r *commentRepository) GetByID(ctx context.Context, id int64) (*entity.Comment, error) {
	const query = `SELECT id, post_id, content, author_id, created_at FROM comments WHERE id = $1`

	var c entity.Comment
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.PostID, &c.Content, &c.AuthorID, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	return &c, nil
}

func (r *commentRepository) ListByPostID(ctx context.Context, postID int64, limit, offset int) ([]*entity.Comment, error) {
	const query = `SELECT id, post_id, content, author_id, created_at FROM comments WHERE post_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		var c entity.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.Content, &c.AuthorID, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return comments, nil
}
