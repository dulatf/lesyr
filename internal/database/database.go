package database

import (
	"context"
	"fmt"
	"time"

	"github.com/dulatf/lesyr/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewDBPool creates a new database connection pool
func NewDBPool(cfg *config.Config) (*pgxpool.Pool, error) {
	// Create a context with timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure the connection pool
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	// Set some reasonable pool settings
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return pool, nil
}
