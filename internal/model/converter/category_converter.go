package converter

import (
	"github.com/tnnz20/jgd-task-1/internal/entity"
	"github.com/tnnz20/jgd-task-1/internal/model"
)

// CategoryToResponse converts entity.Category to model.CategoryResponse
func CategoryToResponse(category *entity.Category) *model.CategoryResponse {
	return &model.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.UnixMilli(),
		UpdatedAt:   category.UpdatedAt.UnixMilli(),
	}
}

// CategoriesToResponses converts slice of entity.Category to slice of model.CategoryResponse
func CategoriesToResponses(categories []*entity.Category) []*model.CategoryResponse {
	responses := make([]*model.CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = CategoryToResponse(category)
	}
	return responses
}
