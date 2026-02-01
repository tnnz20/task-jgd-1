package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/tnnz20/jgd-task-1/internal/model"
	"github.com/tnnz20/jgd-task-1/internal/usecase"
)

// ProductController handles HTTP requests for products
type ProductController struct {
	UseCase *usecase.ProductUseCase
	Log     *slog.Logger
}

// NewProductController creates a new product controller
func NewProductController(useCase *usecase.ProductUseCase, logger *slog.Logger) *ProductController {
	return &ProductController{
		UseCase: useCase,
		Log:     logger,
	}
}

// Create handles POST /api/products
func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	request := new(model.CreateProductRequest)
	if err := ReadJSON(r, request); err != nil {
		c.Log.Warn("Invalid request body", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := c.UseCase.Create(request)
	if err != nil {
		if errors.Is(err, usecase.ErrProductBadRequest) {
			WriteError(w, http.StatusBadRequest, "Invalid product data")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	WriteJSON(w, http.StatusCreated, model.WebResponse[*model.ProductResponse]{Data: response})
}

// List handles GET /api/products
func (c *ProductController) List(w http.ResponseWriter, r *http.Request) {
	responses, err := c.UseCase.List()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	WriteJSON(w, http.StatusOK, model.WebResponse[[]*model.ProductResponse]{Data: responses})
}

// Get handles GET /api/products/{id}
func (c *ProductController) Get(w http.ResponseWriter, r *http.Request) {
	id, err := GetIDFromPath(r, "id")
	if err != nil {
		c.Log.Warn("Invalid product ID", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	request := &model.GetProductRequest{ID: id}
	response, err := c.UseCase.Get(request)
	if err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			WriteError(w, http.StatusNotFound, "Product not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve product")
		return
	}

	WriteJSON(w, http.StatusOK, model.WebResponse[*model.ProductResponse]{Data: response})
}

// Update handles PUT /api/products/{id}
func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := GetIDFromPath(r, "id")
	if err != nil {
		c.Log.Warn("Invalid product ID", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	request := new(model.UpdateProductRequest)
	if err := ReadJSON(r, request); err != nil {
		c.Log.Warn("Invalid request body", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	request.ID = id

	response, err := c.UseCase.Update(request)
	if err != nil {
		if errors.Is(err, usecase.ErrProductBadRequest) {
			WriteError(w, http.StatusBadRequest, "Invalid product data")
			return
		}
		if errors.Is(err, usecase.ErrProductNotFound) {
			WriteError(w, http.StatusNotFound, "Product not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	WriteJSON(w, http.StatusOK, model.WebResponse[*model.ProductResponse]{Data: response})
}

// Delete handles DELETE /api/products/{id}
func (c *ProductController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := GetIDFromPath(r, "id")
	if err != nil {
		c.Log.Warn("Invalid product ID", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	request := &model.DeleteProductRequest{ID: id}
	err = c.UseCase.Delete(request)
	if err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			WriteError(w, http.StatusNotFound, "Product not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	WriteJSON(w, http.StatusOK, model.WebResponse[*model.ProductResponse]{Data: nil})
}
