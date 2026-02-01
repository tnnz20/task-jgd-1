package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tnnz20/jgd-task-1/internal/config"
)

func main() {
	// Initialize Viper configuration
	v := config.NewViper()

	// Initialize logger from viper config
	logger := config.NewLogger(v)

	// Load application config from viper
	appConfig := config.NewConfig(v)

	// Validate required environment variables
	if appConfig.App.Port == "" {
		appConfig.App.Port = "8080"
	}
	if appConfig.App.LogLevel == "" {
		appConfig.App.LogLevel = "INFO"
	}

	logger.Info("Server initializing",
		slog.String("environment", appConfig.App.Environment),
		slog.String("port", appConfig.App.Port),
		slog.String("log_level", appConfig.App.LogLevel),
	)

	// Initialize database connection (optional - only if DB_HOST is provided)
	var db *pgxpool.Pool
	if appConfig.Database.Host != "" {
		db = config.NewDatabase(v, logger)
		defer config.CloseDatabase(db, logger)
	} else {
		logger.Warn("Database not configured, using in-memory repository")
	}

	// Create new ServeMux
	app := http.NewServeMux()

	// Bootstrap application (dependency injection)
	config.Bootstrap(&config.BootstrapConfig{
		App:    app,
		Logger: logger,
		Config: v,
		DB:     db,
	})

	// Create server with configuration
	server := &http.Server{
		Addr:         ":" + appConfig.App.Port,
		Handler:      app,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for errors from server
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		logger.Info("Server starting",
			slog.String("addr", "http://localhost:"+appConfig.App.Port),
		)

		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for interrupt/terminate signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server error", slog.String("error", err.Error()))
		}

	case sig := <-shutdown:
		logger.Info("Shutdown signal received", slog.String("signal", sig.String()))

		// Create context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		logger.Info("Shutting down server gracefully...")

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Graceful shutdown failed, forcing close",
				slog.String("error", err.Error()),
			)
			if err := server.Close(); err != nil {
				logger.Error("Server close error", slog.String("error", err.Error()))
			}
		}

		logger.Info("Server shutdown complete")
	}
}
