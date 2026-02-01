package config

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	deliveryhttp "github.com/tnnz20/jgd-task-1/internal/delivery/http"
	"github.com/tnnz20/jgd-task-1/internal/delivery/http/route"
	"github.com/tnnz20/jgd-task-1/internal/repository"
	"github.com/tnnz20/jgd-task-1/internal/repository/memory"
	"github.com/tnnz20/jgd-task-1/internal/repository/postgres"
	"github.com/tnnz20/jgd-task-1/internal/usecase"
)

// BootstrapConfig holds the configuration for bootstrapping the application
type BootstrapConfig struct {
	App    *http.ServeMux
	Logger *slog.Logger
	Config *viper.Viper
	DB     *pgxpool.Pool
}

// Bootstrap initializes all dependencies and configures routes
func Bootstrap(config *BootstrapConfig) {
	// Setup repositories based on available database
	var categoryRepo repository.CategoryRepositoryInterface

	if config.DB != nil {
		// Use PostgreSQL repository
		config.Logger.Info("Using PostgreSQL repository")
		categoryRepo = postgres.NewCategoryRepository(config.DB)
	} else {
		// Use in-memory repository
		config.Logger.Info("Using in-memory repository")
		categoryRepo = memory.NewCategoryRepository()
	}

	// Setup use cases
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo, config.Logger)

	// Setup controllers
	categoryController := deliveryhttp.NewCategoryController(categoryUseCase, config.Logger)

	// Setup routes
	routeConfig := route.RouteConfig{
		App:                config.App,
		CategoryController: categoryController,
	}
	routeConfig.Setup()
}
