package config

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

// NewLogger creates a new structured logger using slog with config from viper
func NewLogger(v *viper.Viper) *slog.Logger {
	logLevel := v.GetString("LOG_LEVEL")

	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create handler options
	opts := &slog.HandlerOptions{
		Level: level,
	}

	// Create JSON handler for structured logging
	handler := slog.NewJSONHandler(os.Stdout, opts)

	// Create and return logger
	logger := slog.New(handler)

	// Set as default logger
	slog.SetDefault(logger)

	return logger
}
