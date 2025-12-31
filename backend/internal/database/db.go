package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/tau-tau-run/backend/config"
)

var DB *sql.DB

// Connect establishes a connection to PostgreSQL database
func Connect(cfg *config.Config) error {
	var err error
	
	// Open database connection
	DB, err = sql.Open("postgres", cfg.DatabaseDSN())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(cfg.Database.MaxConnections)
	DB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	DB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// HealthCheck checks if database connection is alive
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	return DB.Ping()
}
