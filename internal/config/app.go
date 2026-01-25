package config

import (
	"log/slog"
	"net/http"

	deliveryhttp "github.com/tnnz20/jgd-task-1/internal/delivery/http"
	"github.com/tnnz20/jgd-task-1/internal/delivery/http/route"
	"github.com/tnnz20/jgd-task-1/internal/repository"
	"github.com/tnnz20/jgd-task-1/internal/usecase"
)

// BootstrapConfig holds the configuration for bootstrapping the application
type BootstrapConfig struct {
	App    *http.ServeMux
	Logger *slog.Logger
}

// Bootstrap initializes all dependencies and configures routes
func Bootstrap(config *BootstrapConfig) {
	// Setup repositories
	categoryRepository := repository.NewCategoryRepository()

	// Setup use cases
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository, config.Logger)

	// Setup controllers
	categoryController := deliveryhttp.NewCategoryController(categoryUseCase, config.Logger)

	// Setup routes
	routeConfig := route.RouteConfig{
		App:                config.App,
		CategoryController: categoryController,
	}
	routeConfig.Setup()
}
