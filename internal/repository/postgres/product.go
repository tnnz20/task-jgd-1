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
	ErrProductNotFound = errors.New("product not found")
)

// ProductRepository handles data operations for products using PostgreSQL
type ProductRepository struct {
	pool *pgxpool.Pool
}

// NewProductRepository creates a new PostgreSQL product repository
func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		pool: pool,
	}
}

// Create adds a new product to the database with category join
func (r *ProductRepository) Create(product *entity.Product) error {
	query := `
		INSERT INTO products (name, price, stock, category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(
		context.Background(),
		query,
		product.Name,
		product.Price,
		product.Stock,
		product.CategoryID,
		time.Now(),
		time.Now(),
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Update modifies an existing product in the database
func (r *ProductRepository) Update(product *entity.Product) error {
	// First check if product exists
	var exists bool
	err := r.pool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)",
		product.ID,
	).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return ErrProductNotFound
	}

	query := `
		UPDATE products
		SET name = $1, price = $2, stock = $3, category_id = $4, updated_at = $5
		WHERE id = $6
		RETURNING created_at, updated_at
	`

	err = r.pool.QueryRow(
		context.Background(),
		query,
		product.Name,
		product.Price,
		product.Stock,
		product.CategoryID,
		time.Now(),
		product.ID,
	).Scan(&product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Delete removes a product from the database
func (r *ProductRepository) Delete(product *entity.Product) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.pool.Exec(context.Background(), query, product.ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}

// FindById finds a product by its ID with category information
func (r *ProductRepository) FindById(product *entity.Product, id int) error {
	query := `
		SELECT 
			p.id, p.name, p.price, p.stock, p.category_id, 
			c.name as category_name, 
			p.created_at, p.updated_at
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	err := r.pool.QueryRow(context.Background(), query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.CategoryID,
		&product.CategoryName,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}

	return nil
}

// FindAll returns all products from the database with category information
func (r *ProductRepository) FindAll() ([]*entity.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.stock, p.category_id, 
			c.name as category_name, 
			p.created_at, p.updated_at
		FROM products p
		JOIN categories c ON p.category_id = c.id
		ORDER BY p.id ASC
	`

	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*entity.Product, 0)

	for rows.Next() {
		product := &entity.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Stock,
			&product.CategoryID,
			&product.CategoryName,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// CountById counts products by ID (used for checking existence)
func (r *ProductRepository) CountById(id int) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM products WHERE id = $1`

	err := r.pool.QueryRow(context.Background(), query, id).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
