package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/drewjst/deltagov/internal/models"
)

// Config holds database connection configuration.
type Config struct {
	// URL is the PostgreSQL connection string
	URL string

	// MaxOpenConns sets the maximum number of open connections
	MaxOpenConns int

	// MaxIdleConns sets the maximum number of idle connections
	MaxIdleConns int

	// ConnMaxLifetime sets the maximum lifetime of a connection
	ConnMaxLifetime time.Duration

	// LogLevel sets the GORM logger level
	LogLevel logger.LogLevel
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig(url string) *Config {
	return &Config{
		URL:             url,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		LogLevel:        logger.Warn,
	}
}

// Connect establishes a connection to the PostgreSQL database.
// Returns a configured GORM DB instance.
func Connect(cfg *Config) (*gorm.DB, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("database: DATABASE_URL is required")
	}

	db, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("database: failed to connect: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database: failed to get underlying DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

// Migrate runs auto-migration for all models and creates custom indexes.
func Migrate(db *gorm.DB) error {
	// Run GORM auto-migration
	if err := db.AutoMigrate(
		&models.Bill{},
		&models.Version{},
		&models.Delta{},
	); err != nil {
		return fmt.Errorf("database: auto-migration failed: %w", err)
	}

	// Create GIN index on bills.metadata JSONB column for fast querying
	// Using IF NOT EXISTS to make it idempotent
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_bills_metadata_gin
		ON bills USING GIN (metadata jsonb_path_ops)
	`).Error; err != nil {
		return fmt.Errorf("database: failed to create GIN index on metadata: %w", err)
	}

	// Create GIN index on deltas.delta_json for querying diff data
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_deltas_delta_json_gin
		ON deltas USING GIN (delta_json jsonb_path_ops)
	`).Error; err != nil {
		return fmt.Errorf("database: failed to create GIN index on delta_json: %w", err)
	}

	return nil
}

// Close closes the database connection.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
