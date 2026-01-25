package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/tnnz20/jgd-task-1/internal/model"
	"github.com/tnnz20/jgd-task-1/internal/usecase"
)

// CategoryController handles HTTP requests for categories
type CategoryController struct {
	UseCase *usecase.CategoryUseCase
	Log     *slog.Logger
}

// NewCategoryController creates a new category controller
func NewCategoryController(useCase *usecase.CategoryUseCase, logger *slog.Logger) *CategoryController {
	return &CategoryController{
		UseCase: useCase,
		Log:     logger,
	}
}

// Create handles POST /api/categories
func (c *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	request := new(model.CreateCategoryRequest)
	if err := ReadJSON(r, request); err != nil {
		c.Log.Warn("Invalid request body", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := c.UseCase.Create(request)
	if err != nil {
		if errors.Is(err, usecase.ErrBadRequest) {
			WriteError(w, http.StatusBadRequest, "Name is required")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	WriteJSON(w, http.StatusCreated, model.WebResponse[*model.CategoryResponse]{Data: response})
}

// List handles GET /api/categories
func (c *CategoryController) List(w http.ResponseWriter, r *http.Request) {
	responses, err := c.UseCase.List()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve categories")
		return
	}

	WriteJSON(w, http.StatusOK, model.WebResponse[[]*model.CategoryResponse]{Data: responses})
}

// Get handles GET /api/categories/{id}
func (c *CategoryController) Get(w http.ResponseWriter, r *http.Request) {
	id, err := GetIDFromPath(r, "id")
	if err != nil {
		c.Log.Warn("Invalid category ID", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	request := &model.GetCategoryRequest{ID: id}
	response, err := c.UseCase.Get(request)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "Category not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve category")
		return
	}

	WriteJSON(w, http.StatusOK, model.WebResponse[*model.CategoryResponse]{Data: response})
}

// Update handles PUT /api/categories/{id}
func (c *CategoryController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := GetIDFromPath(r, "id")
	if err != nil {
		c.Log.Warn("Invalid category ID", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	request := new(model.UpdateCategoryRequest)
	if err := ReadJSON(r, request); err != nil {
		c.Log.Warn("Invalid request body", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	request.ID = id

	response, err := c.UseCase.Update(request)
	if err != nil {
		if errors.Is(err, usecase.ErrBadRequest) {
			WriteError(w, http.StatusBadRequest, "Name is required")
			return
		}
		if errors.Is(err, usecase.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "Category not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to update category")
		return
	}

	WriteJSON(w, http.StatusOK, model.WebResponse[*model.CategoryResponse]{Data: response})
}

// Delete handles DELETE /api/categories/{id}
func (c *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := GetIDFromPath(r, "id")
	if err != nil {
		c.Log.Warn("Invalid category ID", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	request := &model.DeleteCategoryRequest{ID: id}
	if err := c.UseCase.Delete(request); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "Category not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	WriteJSON(w, http.StatusOK, model.WebResponse[bool]{Data: true})
}
