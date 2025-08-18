package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/VaneZ444/forum-service/internal/entity"
)

type TagRepo struct {
	db *sql.DB
}

func NewTagRepo(db *sql.DB) *TagRepo {
	return &TagRepo{db: db}
}

func (r *TagRepo) GetByID(ctx context.Context, id int64) (*entity.Tag, error) {
	tag := &entity.Tag{}
	err := r.db.QueryRowContext(ctx, "SELECT id, title, slug FROM tags WHERE id = $1", id).
		Scan(&tag.ID, &tag.Name, &tag.Slug)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (r *TagRepo) GetBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	tag := &entity.Tag{}
	err := r.db.QueryRowContext(ctx, "SELECT id, title, slug FROM tags WHERE slug = $1", slug).
		Scan(&tag.ID, &tag.Name, &tag.Slug)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (r *TagRepo) List(ctx context.Context, limit, offset int) ([]*entity.Tag, int64, error) {
	const q = `SELECT id, name, slug FROM tags ORDER BY name LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tags: %w", err)
	}
	defer rows.Close()

	var tags []*entity.Tag
	for rows.Next() {
		t := new(entity.Tag)
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, 0, err
		}
		tags = append(tags, t)
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tags`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count tags: %w", err)
	}
	return tags, total, nil
}

func (r *TagRepo) ListByPostID(ctx context.Context, postID int64) ([]*entity.Tag, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, t.title, t.slug
		FROM tags t
		JOIN post_tags pt ON pt.tag_id = t.id
		WHERE pt.post_id = $1`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*entity.Tag
	for rows.Next() {
		var t entity.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, err
		}
		tags = append(tags, &t)
	}
	return tags, nil
}

func (r *TagRepo) ListByIDs(ctx context.Context, ids []int64) ([]*entity.Tag, error) {
	if len(ids) == 0 {
		return []*entity.Tag{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf("SELECT id, title, slug FROM tags WHERE id IN (%s)", strings.Join(placeholders, ","))
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*entity.Tag
	for rows.Next() {
		tag := &entity.Tag{}
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (r *TagRepo) Create(ctx context.Context, tag *entity.Tag) (int64, error) {
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO tags (name, slug) VALUES ($1, $2) RETURNING id",
		tag.Name, tag.Slug).Scan(&tag.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to create tag: %w", err)
	}
	return tag.ID, nil
}

func (r *TagRepo) ListAll(ctx context.Context) ([]*entity.Tag, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, slug FROM tags")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*entity.Tag
	for rows.Next() {
		tag := &entity.Tag{}
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
func (r *TagRepo) AddToPost(ctx context.Context, postID int64, tagID int64) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO post_tags (post_id, tag_id) 
        VALUES ($1, $2)`,
		postID, tagID)
	return err
}

func (r *TagRepo) RemoveFromPost(ctx context.Context, postID int64, tagID int64) error {
	_, err := r.db.ExecContext(ctx, `
        DELETE FROM post_tags 
        WHERE post_id = $1 AND tag_id = $2`,
		postID, tagID)
	return err
}
