package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/tnnz20/jgd-task-1/internal/entity"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

// CategoryRepository handles data operations for categories
type CategoryRepository struct {
	mu         sync.RWMutex
	categories []*entity.Category // in-memory storage (empty array)
	counter    int                // auto-increment ID
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		categories: make([]*entity.Category, 0), // initialize empty slice
		counter:    0,
	}
}

// Create adds a new category to the repository
func (r *CategoryRepository) Create(category *entity.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counter++
	category.ID = r.counter
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	r.categories = append(r.categories, category)
	return nil
}

// Update modifies an existing category
func (r *CategoryRepository) Update(category *entity.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, existing := range r.categories {
		if existing.ID == category.ID {
			category.CreatedAt = existing.CreatedAt
			category.UpdatedAt = time.Now()
			r.categories[i] = category
			return nil
		}
	}
	return ErrCategoryNotFound
}

// Delete removes a category from the repository
func (r *CategoryRepository) Delete(category *entity.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, existing := range r.categories {
		if existing.ID == category.ID {
			// Remove by replacing with last element and truncating
			r.categories[i] = r.categories[len(r.categories)-1]
			r.categories = r.categories[:len(r.categories)-1]
			return nil
		}
	}
	return ErrCategoryNotFound
}

// FindById finds a category by its ID
func (r *CategoryRepository) FindById(category *entity.Category, id int) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, existing := range r.categories {
		if existing.ID == id {
			*category = *existing
			return nil
		}
	}
	return ErrCategoryNotFound
}

// FindAll returns all categories
func (r *CategoryRepository) FindAll() ([]*entity.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent external modifications
	result := make([]*entity.Category, len(r.categories))
	copy(result, r.categories)
	return result, nil
}

// CountById counts categories by ID (used for checking existence)
func (r *CategoryRepository) CountById(id int) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, category := range r.categories {
		if category.ID == id {
			return 1, nil
		}
	}
	return 0, nil
}
