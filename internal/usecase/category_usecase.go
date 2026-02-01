package usecase

import (
	"errors"
	"log/slog"

	"github.com/tnnz20/jgd-task-1/internal/entity"
	"github.com/tnnz20/jgd-task-1/internal/model"
	"github.com/tnnz20/jgd-task-1/internal/model/converter"
	"github.com/tnnz20/jgd-task-1/internal/repository"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
	ErrInternal   = errors.New("internal server error")
)

// CategoryUseCase handles business logic for categories
type CategoryUseCase struct {
	CategoryRepository repository.CategoryRepositoryInterface
	Log                *slog.Logger
}

// NewCategoryUseCase creates a new category use case
func NewCategoryUseCase(categoryRepository repository.CategoryRepositoryInterface, logger *slog.Logger) *CategoryUseCase {
	return &CategoryUseCase{
		CategoryRepository: categoryRepository,
		Log:                logger,
	}
}

// Create creates a new category
func (c *CategoryUseCase) Create(request *model.CreateCategoryRequest) (*model.CategoryResponse, error) {
	// Validation
	if request.Name == "" {
		c.Log.Warn("Create category failed: name is required")
		return nil, ErrBadRequest
	}

	category := &entity.Category{
		Name:        request.Name,
		Description: request.Description,
	}

	if err := c.CategoryRepository.Create(category); err != nil {
		c.Log.Error("Failed to create category", slog.String("error", err.Error()))
		return nil, ErrInternal
	}

	c.Log.Info("Category created", slog.Int("id", category.ID), slog.String("name", category.Name))
	return converter.CategoryToResponse(category), nil
}

// Update updates an existing category
func (c *CategoryUseCase) Update(request *model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	// Validation
	if request.Name == "" {
		c.Log.Warn("Update category failed: name is required", slog.Int("id", request.ID))
		return nil, ErrBadRequest
	}

	// Check if category exists
	category := new(entity.Category)
	if err := c.CategoryRepository.FindById(category, request.ID); err != nil {
		c.Log.Warn("Category not found", slog.Int("id", request.ID))
		return nil, ErrNotFound
	}

	// Update category
	category.Name = request.Name
	category.Description = request.Description

	if err := c.CategoryRepository.Update(category); err != nil {
		c.Log.Error("Failed to update category", slog.Int("id", request.ID), slog.String("error", err.Error()))
		return nil, ErrInternal
	}

	c.Log.Info("Category updated", slog.Int("id", category.ID), slog.String("name", category.Name))
	return converter.CategoryToResponse(category), nil
}

// Delete deletes a category
func (c *CategoryUseCase) Delete(request *model.DeleteCategoryRequest) error {
	// Check if category exists
	category := new(entity.Category)
	if err := c.CategoryRepository.FindById(category, request.ID); err != nil {
		c.Log.Warn("Category not found for deletion", slog.Int("id", request.ID))
		return ErrNotFound
	}

	if err := c.CategoryRepository.Delete(category); err != nil {
		c.Log.Error("Failed to delete category", slog.Int("id", request.ID), slog.String("error", err.Error()))
		return ErrInternal
	}

	c.Log.Info("Category deleted", slog.Int("id", request.ID))
	return nil
}

// Get retrieves a category by ID
func (c *CategoryUseCase) Get(request *model.GetCategoryRequest) (*model.CategoryResponse, error) {
	category := new(entity.Category)
	if err := c.CategoryRepository.FindById(category, request.ID); err != nil {
		c.Log.Warn("Category not found", slog.Int("id", request.ID))
		return nil, ErrNotFound
	}

	c.Log.Debug("Category retrieved", slog.Int("id", category.ID))
	return converter.CategoryToResponse(category), nil
}

// List retrieves all categories
func (c *CategoryUseCase) List() ([]*model.CategoryResponse, error) {
	categories, err := c.CategoryRepository.FindAll()
	if err != nil {
		c.Log.Error("Failed to list categories", slog.String("error", err.Error()))
		return nil, ErrInternal
	}

	c.Log.Debug("Categories listed", slog.Int("count", len(categories)))
	return converter.CategoriesToResponses(categories), nil
}
