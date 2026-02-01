package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	ErrFailedToGetNextVersion = "failed to get next version %w"
	ErrCreateMigrationDriver  = "could not create migration driver: %w"
	ErrCreateMigrateInstance  = "could not create migrate instance: %w"
)

// getNextVersion scans the migrations directory to find the next version number
func getNextVersion(migrationsDir string) (int, error) {
	// Read all files in migrations directory
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return 0, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	maxVersion := 0

	// Find the highest version number
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Migration files start with 6-digit version number
		if len(name) >= 6 {
			var version int
			_, err := fmt.Sscanf(name[:6], "%d", &version)
			if err == nil && version > maxVersion {
				maxVersion = version
			}
		}
	}

	return maxVersion + 1, nil
}

// CreateMigrationFiles creates new up and down migration files with the given name
func CreateMigrationFiles(name, migrationsDir string) error {

	// Find the next version number
	nextVersion, err := getNextVersion(migrationsDir)
	if err != nil {
		return fmt.Errorf(ErrFailedToGetNextVersion, err)
	}

	// Create file names
	upFile := fmt.Sprintf("%s/%06d_%s.up.sql", migrationsDir, nextVersion, name)
	downFile := fmt.Sprintf("%s/%06d_%s.down.sql", migrationsDir, nextVersion, name)

	// Create UP migration file
	upContent := fmt.Sprintf("-- Migration: %s\n-- Created: %s\n\n-- Write your UP migration here\n",
		name, time.Now().Format("2006-01-02 15:04:05"))
	if err := os.WriteFile(upFile, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create UP file: %w", err)
	}

	// Create DOWN migration file
	downContent := fmt.Sprintf("-- Migration: %s\n-- Created: %s\n\n-- Write your DOWN migration here\n",
		name, time.Now().Format("2006-01-02 15:04:05"))
	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to create DOWN file: %w", err)
	}

	log.Printf("✅ Created migration files:")
	log.Printf("   UP:   %s", upFile)
	log.Printf("   DOWN: %s", downFile)

	return nil
}

func RunMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf(ErrCreateMigrationDriver, err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf(ErrCreateMigrateInstance, err)
	}

	// Run all pending migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}

	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf(ErrCreateMigrationDriver, err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf(ErrCreateMigrateInstance, err)
	}

	// Rollback one migration
	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not rollback migration: %w", err)
	}

	log.Println("Migration rolled back successfully")
	return nil
}
