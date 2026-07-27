package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// creates a postgres pool with retry logic
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	const maxRetries = 10
	const retryDelay = 3 * time.Second

	var pool *pgxpool.Pool
	var err error

	for i := range maxRetries {
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			slog.Warn("failed to create pool, retrying...",
				"attempt", i+1,
				"error", err,
			)
			time.Sleep(retryDelay)
			continue
		}

		if err = pool.Ping(ctx); err != nil {
			pool.Close()
			slog.Warn("failed to ping database, retrying...",
				"attempt", i+1,
				"error", err,
			)
			time.Sleep(retryDelay)
			continue
		}

		slog.Info("connected to PostgreSQL")
		return pool, nil
	}

	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d attempts: %w", maxRetries, err)
}
