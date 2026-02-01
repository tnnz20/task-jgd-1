package route

import (
	"net/http"

	deliveryhttp "github.com/tnnz20/jgd-task-1/internal/delivery/http"
)

// RouteConfig holds the configuration for routes
type RouteConfig struct {
	App                *http.ServeMux
	CategoryController *deliveryhttp.CategoryController
	ProductController  *deliveryhttp.ProductController
}

// Setup configures all routes
func (c *RouteConfig) Setup() {
	c.SetupCategoryRoute()
	c.SetupProductRoute()
}

// SetupCategoryRoute configures category routes
func (c *RouteConfig) SetupCategoryRoute() {
	c.App.HandleFunc("POST /api/categories", c.CategoryController.Create)
	c.App.HandleFunc("GET /api/categories", c.CategoryController.List)
	c.App.HandleFunc("GET /api/categories/{id}", c.CategoryController.Get)
	c.App.HandleFunc("PUT /api/categories/{id}", c.CategoryController.Update)
	c.App.HandleFunc("DELETE /api/categories/{id}", c.CategoryController.Delete)

	// Health check endpoint
	c.App.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})
}

// SetupProductRoute configures product routes
func (c *RouteConfig) SetupProductRoute() {
	c.App.HandleFunc("POST /api/products", c.ProductController.Create)
	c.App.HandleFunc("GET /api/products", c.ProductController.List)
	c.App.HandleFunc("GET /api/products/{id}", c.ProductController.Get)
	c.App.HandleFunc("PUT /api/products/{id}", c.ProductController.Update)
	c.App.HandleFunc("DELETE /api/products/{id}", c.ProductController.Delete)
}
