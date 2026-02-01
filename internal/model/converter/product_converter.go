package converter

import (
	"github.com/tnnz20/jgd-task-1/internal/entity"
	"github.com/tnnz20/jgd-task-1/internal/model"
)

// ProductToResponse converts entity.Product to model.ProductResponse
func ProductToResponse(product *entity.Product) *model.ProductResponse {
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

// ProductsToResponses converts slice of entity.Product to slice of model.ProductResponse
func ProductsToResponses(products []*entity.Product) []*model.ProductResponse {
	responses := make([]*model.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = ProductToResponse(product)
	}
	return responses
}
