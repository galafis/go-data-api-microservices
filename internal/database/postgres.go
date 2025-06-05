package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/galafis/go-data-api-microservices/internal/config"
	"github.com/galafis/go-data-api-microservices/pkg/logger"
	_ "github.com/lib/pq"
)

// PostgresDB represents a PostgreSQL database connection
type PostgresDB struct {
	DB *sql.DB
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg *config.DatabaseConfig) (*PostgresDB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Connected to PostgreSQL database")
	return &PostgresDB{DB: db}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	if p.DB != nil {
		return p.DB.Close()
	}
	return nil
}

// Ping checks if the database connection is alive
func (p *PostgresDB) Ping() error {
	return p.DB.Ping()
}

// Begin starts a new transaction
func (p *PostgresDB) Begin() (*sql.Tx, error) {
	return p.DB.Begin()
}

// Query executes a query that returns rows
func (p *PostgresDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.DB.Query(query, args...)
}

// QueryRow executes a query that returns a single row
func (p *PostgresDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return p.DB.QueryRow(query, args...)
}

// Exec executes a query that doesn't return rows
func (p *PostgresDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.DB.Exec(query, args...)
}

// PrepareStatement prepares a statement for later queries or executions
func (p *PostgresDB) PrepareStatement(query string) (*sql.Stmt, error) {
	return p.DB.Prepare(query)
}

// ExecuteWithRetry executes a function with retry logic
func (p *PostgresDB) ExecuteWithRetry(fn func() error, maxRetries int, delay time.Duration) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		
		logger.Warnf("Database operation failed (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(delay)
			delay *= 2 // Exponential backoff
		}
	}
	return err
}

