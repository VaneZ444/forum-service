package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/VaneZ444/forum-service/internal/entity"
	"github.com/VaneZ444/forum-service/internal/repository"
)

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) repository.CategoryRepository {
	return &categoryRepository{db: db}
}
func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) error {
	query := `UPDATE categories 
              SET title = $1, description = $2, updated_at = $3 
              WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query,
		category.Title,
		category.Description,
		category.UpdatedAt,
		category.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found with id: %d", category.ID)
	}

	return nil
}
func (r *categoryRepository) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	const query = `SELECT id, title, slug FROM categories WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var category entity.Category
	if err := row.Scan(&category.ID, &category.Title, &category.Slug); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get category by ID: %w", err)
	}

	return &category, nil
}

func (r *categoryRepository) List(ctx context.Context, limit, offset int) ([]*entity.Category, error) {
	query := `SELECT id, title, slug, description, created_at, updated_at 
              FROM categories 
              ORDER BY created_at DESC
              LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		var c entity.Category
		if err := rows.Scan(&c.ID, &c.Title, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) (int64, error) {
	const query = `INSERT INTO categories (title, slug) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, category.Title, category.Slug).Scan(&category.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to create category: %w", err)
	}

	return category.ID, nil
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	const query = `SELECT id, title, slug FROM categories WHERE slug = $1`

	row := r.db.QueryRowContext(ctx, query, slug)

	var category entity.Category
	if err := row.Scan(&category.ID, &category.Title, &category.Slug); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found by slug: %w", err)
		}
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	return &category, nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category not found with id: %d", id)
	}
	return nil
}
