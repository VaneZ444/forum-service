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

func (r *commentRepository) ListByPost(ctx context.Context, postID int64, limit, offset int) ([]*entity.Comment, int64, error) {
	const countQ = `SELECT COUNT(*) FROM comments WHERE post_id = $1`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQ, postID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	const q = `SELECT id, post_id, content, author_id, created_at
               FROM comments WHERE post_id = $1
               ORDER BY created_at ASC
               LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, q, postID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list comments: %w", err)
	}
	defer rows.Close()

	var items []*entity.Comment
	for rows.Next() {
		c := new(entity.Comment)
		if err := rows.Scan(&c.ID, &c.PostID, &c.Content, &c.AuthorID, &c.CreatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, c)
	}
	return items, total, nil
}
