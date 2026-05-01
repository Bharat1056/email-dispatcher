package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/user/queue/internal/config"
)

func NewPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {

	// parse the connection string to get the database configuration
	dbConfig, err := pgxpool.ParseConfig(cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	// setup DB configurations
	dbConfig.MaxConns = cfg.DBMaxConns
	dbConfig.MinConns = cfg.DBMinConns
	dbConfig.MaxConnLifetime = cfg.DBMaxConnLifetime
	dbConfig.MaxConnIdleTime = cfg.DBMaxConnIdleTime
	dbConfig.HealthCheckPeriod = cfg.DBHealthCheckPeriod

	// create the database pool using the parsed configuration
	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Verify connection via ping the database
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging the database: %w", err)
	}

	log.Printf("Successfully connected to database with %d max connections", cfg.DBMaxConns)
	return pool, nil
}
