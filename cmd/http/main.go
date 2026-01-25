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

	"github.com/tnnz20/jgd-task-1/internal/config"
)

func main() {
	// Initialize logger
	logger := config.NewLogger()

	// Create new ServeMux
	app := http.NewServeMux()

	// Bootstrap application (dependency injection)
	config.Bootstrap(&config.BootstrapConfig{
		App:    app,
		Logger: logger,
	})

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create server with configuration
	server := &http.Server{
		Addr:         ":" + port,
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
			slog.String("addr", "http://localhost:"+port),
		)
		logger.Info("API Endpoints",
			slog.String("health", "GET /health"),
			slog.String("create", "POST /api/categories"),
			slog.String("list", "GET /api/categories"),
			slog.String("get", "GET /api/categories/{id}"),
			slog.String("update", "PUT /api/categories/{id}"),
			slog.String("delete", "DELETE /api/categories/{id}"),
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

		logger.Info("Server stopped")
	}
}
