package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const migrationSQL = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS webhooks (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source      VARCHAR(100) NOT NULL DEFAULT 'unknown',
    payload     JSONB NOT NULL,
    target_url  TEXT NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    retry_count INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_webhooks_status ON webhooks(status);
CREATE INDEX IF NOT EXISTS idx_webhooks_created_at ON webhooks(created_at DESC);
`

// creates a postgres pool with retry logic and runs migrations
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

		// auto-migrate schema
		if _, err := pool.Exec(ctx, migrationSQL); err != nil {
			pool.Close()
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
		slog.Info("database migrations applied")

		return pool, nil
	}

	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d attempts: %w", maxRetries, err)
}
