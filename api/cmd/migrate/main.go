package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/portnd/the-sentinel-core/internal/core/config"
	_ "github.com/lib/pq"
)

const (
	migrationsDir = "databases/migrations"
	upSuffix      = ".up.sql"
	downSuffix    = ".down.sql"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Get command (up or down)
	if len(os.Args) < 2 {
		log.Fatal("❌ Usage: go run cmd/migrate/main.go [up|down]")
	}
	command := os.Args[1]

	// Validate command
	if command != "up" && command != "down" {
		log.Fatalf("❌ Invalid command: %s. Use 'up' or 'down'", command)
	}

	// Connect to PostgreSQL
	log.Printf("🔌 Connecting to PostgreSQL at %s:%s...", cfg.PostgresHost, cfg.PostgresPort)
	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}
	log.Println("✅ Database connection established")

	// Create migrations table if not exists
	if err := createMigrationsTable(db); err != nil {
		log.Fatalf("❌ Failed to create migrations table: %v", err)
	}

	// Run migrations
	if command == "up" {
		if err := runMigrationsUp(db); err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}
		log.Println("✅ All migrations completed successfully!")
	} else {
		if err := runMigrationsDown(db); err != nil {
			log.Fatalf("❌ Rollback failed: %v", err)
		}
		log.Println("✅ Rollback completed successfully!")
	}
}

// createMigrationsTable creates the schema_migrations table to track applied migrations
func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			version VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`
	_, err := db.Exec(query)
	return err
}

// runMigrationsUp applies all pending migrations
func runMigrationsUp(db *sql.DB) error {
	// Get list of migration files
	files, err := getMigrationFiles(upSuffix)
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	if len(files) == 0 {
		log.Println("ℹ️  No migration files found")
		return nil
	}

	// Get already applied migrations
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply pending migrations
	appliedCount := 0
	for _, file := range files {
		version := extractVersion(file)

		// Skip if already applied
		if _, exists := appliedMigrations[version]; exists {
			log.Printf("⏭️  Skipping %s (already applied)", version)
			continue
		}

		// Read migration file
		filePath := filepath.Join(migrationsDir, file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", filePath, err)
		}

		// Execute migration
		log.Printf("🚀 Applying migration: %s", version)
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute %s: %w", version, err)
		}

		// Record migration
		if err := recordMigration(db, version); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", version, err)
		}

		log.Printf("✅ Applied: %s", version)
		appliedCount++
	}

	if appliedCount == 0 {
		log.Println("ℹ️  Database is up to date. No migrations to apply.")
	} else {
		log.Printf("✅ Applied %d migration(s)", appliedCount)
	}

	return nil
}

// runMigrationsDown rolls back the last applied migration
func runMigrationsDown(db *sql.DB) error {
	// Get list of migration files
	files, err := getMigrationFiles(downSuffix)
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	if len(files) == 0 {
		log.Println("ℹ️  No migration files found")
		return nil
	}

	// Get the last applied migration
	lastMigration, err := getLastAppliedMigration(db)
	if err != nil {
		return fmt.Errorf("failed to get last applied migration: %w", err)
	}

	if lastMigration == "" {
		log.Println("ℹ️  No migrations to rollback")
		return nil
	}

	// Find the corresponding down file
	downFile := lastMigration + downSuffix
	filePath := filepath.Join(migrationsDir, downFile)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("rollback file not found: %s", downFile)
	}

	// Read down migration file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	// Execute rollback
	log.Printf("🔄 Rolling back migration: %s", lastMigration)
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute rollback for %s: %w", lastMigration, err)
	}

	// Remove migration record
	if err := removeMigration(db, lastMigration); err != nil {
		return fmt.Errorf("failed to remove migration record %s: %w", lastMigration, err)
	}

	log.Printf("✅ Rolled back: %s", lastMigration)
	return nil
}

// getMigrationFiles returns sorted list of migration files with given suffix
func getMigrationFiles(suffix string) ([]string, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), suffix) {
			files = append(files, entry.Name())
		}
	}

	// Sort files by name (timestamp prefix ensures correct order)
	sort.Strings(files)
	return files, nil
}

// extractVersion extracts the version from filename (removes suffix)
func extractVersion(filename string) string {
	version := strings.TrimSuffix(filename, upSuffix)
	version = strings.TrimSuffix(version, downSuffix)
	return version
}

// getAppliedMigrations returns a map of already applied migrations
func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	migrations := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		migrations[version] = true
	}

	return migrations, rows.Err()
}

// getLastAppliedMigration returns the version of the last applied migration
func getLastAppliedMigration(db *sql.DB) (string, error) {
	var version string
	err := db.QueryRow("SELECT version FROM schema_migrations ORDER BY applied_at DESC LIMIT 1").Scan(&version)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return version, err
}

// recordMigration records a migration as applied
func recordMigration(db *sql.DB, version string) error {
	_, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
	return err
}

// removeMigration removes a migration record
func removeMigration(db *sql.DB, version string) error {
	_, err := db.Exec("DELETE FROM schema_migrations WHERE version = $1", version)
	return err
}
