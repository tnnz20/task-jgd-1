package config

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

// NewDatabase creates and returns a new PostgreSQL connection pool
func NewDatabase(v *viper.Viper, logger *slog.Logger) *pgxpool.Pool {
	dbConfig := &DatabaseConfig{
		Host:     v.GetString("DB_HOST"),
		Port:     v.GetString("DB_PORT"),
		Name:     v.GetString("DB_NAME"),
		User:     v.GetString("DB_USER"),
		Password: v.GetString("DB_PASSWORD"),
		PoolMode: v.GetString("DB_POOLMODE"),
	}

	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_mode=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.PoolMode,
	)

	// Create connection pool
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Error("Failed to create database connection pool",
			slog.String("error", err.Error()),
		)
		panic(fmt.Errorf("unable to create connection pool: %w", err))
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		logger.Error("Failed to ping database",
			slog.String("error", err.Error()),
		)
		panic(fmt.Errorf("unable to ping database: %w", err))
	}

	logger.Info("Database connection established",
		slog.String("host", dbConfig.Host),
		slog.String("port", dbConfig.Port),
		slog.String("name", dbConfig.Name),
	)

	return pool
}

// CloseDatabase closes the database connection pool
func CloseDatabase(pool *pgxpool.Pool, logger *slog.Logger) {
	if pool != nil {
		pool.Close()
		logger.Info("Database connection closed")
	}
}
