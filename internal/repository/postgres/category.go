package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tnnz20/jgd-task-1/internal/entity"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

// CategoryRepository handles data operations for categories using PostgreSQL
type CategoryRepository struct {
	pool *pgxpool.Pool
}

// NewCategoryRepository creates a new PostgreSQL category repository
func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		pool: pool,
	}
}

// Create adds a new category to the database
func (r *CategoryRepository) Create(category *entity.Category) error {
	query := `
		INSERT INTO categories (name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(
		context.Background(),
		query,
		category.Name,
		category.Description,
		time.Now(),
		time.Now(),
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Update modifies an existing category in the database
func (r *CategoryRepository) Update(category *entity.Category) error {
	// First check if category exists
	var exists bool
	err := r.pool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)",
		category.ID,
	).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return ErrCategoryNotFound
	}

	query := `
		UPDATE categories
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
		RETURNING created_at, updated_at
	`

	err = r.pool.QueryRow(
		context.Background(),
		query,
		category.Name,
		category.Description,
		time.Now(),
		category.ID,
	).Scan(&category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Delete removes a category from the database
func (r *CategoryRepository) Delete(category *entity.Category) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.pool.Exec(context.Background(), query, category.ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

// FindById finds a category by its ID
func (r *CategoryRepository) FindById(category *entity.Category, id int) error {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	err := r.pool.QueryRow(context.Background(), query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCategoryNotFound
		}
		return err
	}

	return nil
}

// FindAll returns all categories from the database
func (r *CategoryRepository) FindAll() ([]*entity.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY id ASC
	`

	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*entity.Category, 0)

	for rows.Next() {
		category := &entity.Category{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// CountById counts categories by ID (used for checking existence)
func (r *CategoryRepository) CountById(id int) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM categories WHERE id = $1`

	err := r.pool.QueryRow(context.Background(), query, id).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
