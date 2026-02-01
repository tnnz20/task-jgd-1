package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/tnnz20/jgd-task-1/internal/entity"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

// ProductRepository handles data operations for products in-memory
type ProductRepository struct {
	mu       sync.RWMutex
	products []*entity.Product // in-memory storage
	counter  int               // auto-increment ID
}

// NewProductRepository creates a new in-memory product repository
func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: make([]*entity.Product, 0),
		counter:  0,
	}
}

// Create adds a new product to the repository
func (r *ProductRepository) Create(product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counter++
	product.ID = r.counter
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	r.products = append(r.products, product)
	return nil
}

// Update modifies an existing product
func (r *ProductRepository) Update(product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, existing := range r.products {
		if existing.ID == product.ID {
			product.CreatedAt = existing.CreatedAt
			product.UpdatedAt = time.Now()
			r.products[i] = product
			return nil
		}
	}

	return ErrProductNotFound
}

// Delete removes a product from the repository
func (r *ProductRepository) Delete(product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, existing := range r.products {
		if existing.ID == product.ID {
			r.products = append(r.products[:i], r.products[i+1:]...)
			return nil
		}
	}

	return ErrProductNotFound
}

// FindById retrieves a single product by ID
func (r *ProductRepository) FindById(product *entity.Product, id int) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.products {
		if p.ID == id {
			*product = *p
			return nil
		}
	}

	return ErrProductNotFound
}

// FindAll retrieves all products
func (r *ProductRepository) FindAll() ([]*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.products, nil
}

// CountById checks if a product with the given ID exists
func (r *ProductRepository) CountById(id int) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, product := range r.products {
		if product.ID == id {
			return 1, nil
		}
	}

	return 0, nil
}
