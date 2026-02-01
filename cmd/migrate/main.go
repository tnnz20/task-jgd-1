package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/tnnz20/jgd-task-1/internal/config"
)

type migrationFlags struct {
	create string
	up     bool
	down   bool
}

func main() {

	flags := &migrationFlags{}

	// Define command-line flags
	flag.StringVar(&flags.create, "create", "", "Create a new migration file with the given name")
	flag.BoolVar(&flags.up, "up", false, "Run all pending migrations")
	flag.BoolVar(&flags.down, "down", false, "Rollback the last migration")
	flag.Parse()

	fmt.Println(flags.create)

	migrationsDir := "cmd/migrate/migrations"

	if flags.create != "" {
		if err := CreateMigrationFiles(flags.create, migrationsDir); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}
		return

	}

	migrationsPath := fmt.Sprintf("file://%s", migrationsDir)

	// Initialize Viper configuration
	v := config.NewViper()

	// Initialize logger from viper config
	logger := config.NewLogger(v)

	// Load application config from viper
	appConfig := config.NewConfig(v)

	// Create database connection pool (single reusable connection)
	var db *sql.DB
	if appConfig.Database.Host != "" {
		// Build DSN once
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			appConfig.Database.User,
			appConfig.Database.Password,
			appConfig.Database.Host,
			appConfig.Database.Port,
			appConfig.Database.Name,
		)

		// Open database connection once
		sqlDB, err := sql.Open("pgx", dsn)
		if err != nil {
			logger.Error("Failed to open database", slog.String("error", err.Error()))
			log.Fatalf("Failed to open database: %v", err)
		}
		defer sqlDB.Close()

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := sqlDB.PingContext(ctx); err != nil {
			cancel()
			logger.Error("Failed to ping database", slog.String("error", err.Error()))
			log.Fatalf("Failed to ping database: %v", err)
		}
		cancel()

		db = sqlDB
		logger.Info("Database connection established for migrations")
	} else {
		logger.Warn("Database not configured, skipping migrations")
		return
	}

	executeMigrationCommand(flags, db, migrationsPath)
}

func executeMigrationCommand(flags *migrationFlags, db *sql.DB, migrationsPath string) {
	switch {
	case flags.up:
		if err := RunMigrations(db, migrationsPath); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	case flags.down:
		if err := RollbackMigration(db, migrationsPath); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
	default:
		log.Println("No migration command provided. Use -up to run migrations or -down to rollback the last migration.")
		os.Exit(1)
	}
}
