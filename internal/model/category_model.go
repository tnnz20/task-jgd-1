package model

// CategoryResponse represents the response for category
type CategoryResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// CreateCategoryRequest represents the request for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateCategoryRequest represents the request for updating a category
type UpdateCategoryRequest struct {
	ID          int    `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetCategoryRequest represents the request for getting a category
type GetCategoryRequest struct {
	ID int `json:"-"`
}

// DeleteCategoryRequest represents the request for deleting a category
type DeleteCategoryRequest struct {
	ID int `json:"-"`
}
