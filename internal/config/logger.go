package config

import (
	"log/slog"
	"os"
)

// NewLogger creates a new structured logger using slog (Go 1.21+)
func NewLogger() *slog.Logger {
	// Get log level from environment variable
	level := os.Getenv("LOG_LEVEL")

	var logLevel slog.Level
	switch level {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// Create handler options
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	// Create JSON handler for structured logging
	handler := slog.NewJSONHandler(os.Stdout, opts)

	// Create and return logger
	logger := slog.New(handler)

	// Set as default logger
	slog.SetDefault(logger)

	return logger
}
