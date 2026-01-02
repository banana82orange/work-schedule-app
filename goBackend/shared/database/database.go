package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Pool represents a database connection pool
type Pool struct {
	db   *sql.DB
	once sync.Once
}

var (
	instance *Pool
	once     sync.Once
)

// NewPool creates a new database connection pool
func NewPool(cfg Config) (*Pool, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")
	return &Pool{db: db}, nil
}

// GetDB returns the database connection
func (p *Pool) GetDB() *sql.DB {
	return p.db
}

// Close closes the database connection
func (p *Pool) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// Transaction executes a function within a database transaction
func (p *Pool) Transaction(fn func(*sql.Tx) error) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx failed: %v, rollback failed: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// DefaultConfig returns default database configuration
func DefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "123456789",
		DBName:   "postgres",
		SSLMode:  "disable",
	}
}
