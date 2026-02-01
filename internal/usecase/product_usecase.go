package usecase

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/tnnz20/jgd-task-1/internal/entity"
	"github.com/tnnz20/jgd-task-1/internal/model"
	"github.com/tnnz20/jgd-task-1/internal/repository"
)

var (
	ErrProductBadRequest = errors.New("bad request")
	ErrProductNotFound   = errors.New("product not found")
)

// ProductUseCase handles business logic for products
type ProductUseCase struct {
	ProductRepository repository.ProductRepositoryInterface
	Log               *slog.Logger
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(productRepo repository.ProductRepositoryInterface, log *slog.Logger) *ProductUseCase {
	return &ProductUseCase{
		ProductRepository: productRepo,
		Log:               log,
	}
}

// Create creates a new product
func (u *ProductUseCase) Create(req *model.CreateProductRequest) (*model.ProductResponse, error) {
	// Validation
	if strings.TrimSpace(req.Name) == "" {
		u.Log.Warn("Create product failed: empty name")
		return nil, ErrProductBadRequest
	}

	if req.Price <= 0 {
		u.Log.Warn("Create product failed: invalid price")
		return nil, ErrProductBadRequest
	}

	if req.Stock < 0 {
		u.Log.Warn("Create product failed: invalid stock")
		return nil, ErrProductBadRequest
	}

	if req.CategoryID <= 0 {
		u.Log.Warn("Create product failed: invalid category_id")
		return nil, ErrProductBadRequest
	}

	product := &entity.Product{
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
	}

	err := u.ProductRepository.Create(product)
	if err != nil {
		u.Log.Error("Create product error", slog.String("error", err.Error()))
		return nil, err
	}

	u.Log.Info("Product created", slog.Int("id", product.ID), slog.String("name", product.Name))

	return productToResponse(product), nil
}

// Get retrieves a single product by ID
func (u *ProductUseCase) Get(req *model.GetProductRequest) (*model.ProductResponse, error) {
	product := &entity.Product{}
	err := u.ProductRepository.FindById(product, req.ID)
	if err != nil {
		u.Log.Warn("Get product not found", slog.Int("id", req.ID))
		return nil, createError(ErrProductNotFound)
	}

	return productToResponse(product), nil
}

// List retrieves all products
func (u *ProductUseCase) List() ([]*model.ProductResponse, error) {
	products, err := u.ProductRepository.FindAll()
	if err != nil {
		u.Log.Error("List products error", slog.String("error", err.Error()))
		return nil, err
	}

	responses := make([]*model.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = productToResponse(product)
	}

	return responses, nil
}

// Update updates an existing product
func (u *ProductUseCase) Update(req *model.UpdateProductRequest) (*model.ProductResponse, error) {
	// Validation
	if strings.TrimSpace(req.Name) == "" {
		u.Log.Warn("Update product failed: empty name")
		return nil, createError(ErrProductBadRequest)
	}

	if req.Price <= 0 {
		u.Log.Warn("Update product failed: invalid price")
		return nil, createError(ErrProductBadRequest)
	}

	if req.Stock < 0 {
		u.Log.Warn("Update product failed: invalid stock")
		return nil, createError(ErrProductBadRequest)
	}

	if req.CategoryID <= 0 {
		u.Log.Warn("Update product failed: invalid category_id")
		return nil, createError(ErrProductBadRequest)
	}

	product := &entity.Product{
		ID:         req.ID,
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
	}

	err := u.ProductRepository.Update(product)
	if err != nil {
		u.Log.Warn("Update product not found", slog.Int("id", req.ID))
		if strings.Contains(err.Error(), "not found") {
			return nil, createError(ErrProductNotFound)
		}
		u.Log.Error("Update product error", slog.String("error", err.Error()))
		return nil, err
	}

	u.Log.Info("Product updated", slog.Int("id", product.ID))

	return productToResponse(product), nil
}

// Delete deletes a product
func (u *ProductUseCase) Delete(req *model.DeleteProductRequest) error {
	product := &entity.Product{ID: req.ID}
	err := u.ProductRepository.Delete(product)
	if err != nil {
		u.Log.Warn("Delete product not found", slog.Int("id", req.ID))
		return createError(ErrProductNotFound)
	}

	u.Log.Info("Product deleted", slog.Int("id", req.ID))
	return nil
}

// Helper function to convert entity to response
func productToResponse(product *entity.Product) *model.ProductResponse {
	return &model.ProductResponse{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
		Category: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			ID:   product.CategoryID,
			Name: product.CategoryName,
		},
		CreatedAt: product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Helper function to create errors
func createError(msg error) error {
	return msg
}
