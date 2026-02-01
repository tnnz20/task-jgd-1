package repository

import "github.com/tnnz20/jgd-task-1/internal/entity"

// CategoryRepositoryInterface defines the contract for category repositories
type CategoryRepositoryInterface interface {
	Create(category *entity.Category) error
	Update(category *entity.Category) error
	Delete(category *entity.Category) error
	FindById(category *entity.Category, id int) error
	FindAll() ([]*entity.Category, error)
	CountById(id int) (int64, error)
}

// ProductRepositoryInterface defines the contract for product repositories
type ProductRepositoryInterface interface {
	Create(product *entity.Product) error
	Update(product *entity.Product) error
	Delete(product *entity.Product) error
	FindById(product *entity.Product, id int) error
	FindAll() ([]*entity.Product, error)
	CountById(id int) (int64, error)
}
