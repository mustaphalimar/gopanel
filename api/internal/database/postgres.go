package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mustaphalimar/gopanel/pkg/config"
	"github.com/mustaphalimar/gopanel/pkg/logger"
)

type Postgres struct {
	Pool   *pgxpool.Pool
	logger *logger.Logger
}

func NewPostgres(cfg *config.DatabaseConfig, log *logger.Logger) (*Postgres, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.GetDatabaseURI())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxIdleConns)
	poolConfig.MinConns = 1
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	// create a pool connection
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	log.Info(
		"Database connection established",
		"host", cfg.Host,
		"database", cfg.DBName,
		"max_conns", cfg.MaxOpenConns,
	)

	return &Postgres{
		Pool:   pool,
		logger: log,
	}, nil

}

func (db *Postgres) Close() {
	db.logger.Info("closing database connection pool")
	db.Pool.Close()
}

func (db *Postgres) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

func (db *Postgres) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("database unhealthy: %w", err)
	}

	return nil
}
