package postgres

import (
	"context"
	"database/sql"
	"errors"
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

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	const query = `
        INSERT INTO categories (title, slug, description)
        VALUES ($1, $2, $3)
        RETURNING id, title, slug, description, created_at, updated_at
    `
	newCategory := &entity.Category{}
	err := r.db.QueryRowContext(ctx, query,
		category.Title,
		category.Slug,
		category.Description,
	).Scan(
		&newCategory.ID,
		&newCategory.Title,
		&newCategory.Slug,
		&newCategory.Description,
		&newCategory.CreatedAt,
		&newCategory.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return newCategory, nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	const query = `
		SELECT id, title, slug, description, created_at, updated_at 
		FROM categories 
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	category := &entity.Category{}

	err := row.Scan(
		&category.ID,
		&category.Title,
		&category.Slug,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return category, nil
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	const query = `
		SELECT id, title, slug, description, created_at, updated_at 
		FROM categories 
		WHERE slug = $1
	`

	row := r.db.QueryRowContext(ctx, query, slug)
	category := &entity.Category{}

	err := row.Scan(
		&category.ID,
		&category.Title,
		&category.Slug,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}
	return category, nil
}

func (r *categoryRepository) List(ctx context.Context, limit, offset int) ([]*entity.Category, error) {
	const query = `
		SELECT id, title, slug, description, created_at, updated_at 
		FROM categories 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	categories := []*entity.Category{}
	for rows.Next() {
		var c entity.Category
		if err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Slug,
			&c.Description,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return categories, nil
}

func (r *categoryRepository) Count(ctx context.Context) (int64, error) {
	const query = `SELECT COUNT(*) FROM categories`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count categories: %w", err)
	}
	return count, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	const query = `
		UPDATE categories 
		SET title = $1, slug = $2, description = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, title, slug, description, created_at, updated_at
	`

	updatedCategory := &entity.Category{}
	err := r.db.QueryRowContext(ctx, query,
		category.Title,
		category.Slug,
		category.Description,
		category.UpdatedAt,
		category.ID,
	).Scan(
		&updatedCategory.ID,
		&updatedCategory.Title,
		&updatedCategory.Slug,
		&updatedCategory.Description,
		&updatedCategory.CreatedAt,
		&updatedCategory.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return updatedCategory, nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM categories WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
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
