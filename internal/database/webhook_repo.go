package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// a row in the webhooks table
type Webhook struct {
	ID         uuid.UUID              `json:"id"`
	Source     string                 `json:"source"`
	Payload    map[string]interface{} `json:"payload"`
	TargetURL  string                 `json:"target_url"`
	Status     string                 `json:"status"`
	RetryCount int                    `json:"retry_count"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// crud operations for webhooks
type WebhookRepo struct {
	pool *pgxpool.Pool
}

// creates a new repo
func NewWebhookRepo(pool *pgxpool.Pool) *WebhookRepo {
	return &WebhookRepo{pool: pool}
}

// inserts a webhook and returns its id
func (r *WebhookRepo) Insert(ctx context.Context, source string, payload map[string]interface{}, targetURL string) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO webhooks (id, source, payload, target_url, status, retry_count, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, 'PENDING', 0, NOW(), NOW())`,
		id, source, payload, targetURL,
	)
	return id, err
}

// updates status and retry count
func (r *WebhookRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, retryCount int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE webhooks SET status = $1, retry_count = $2, updated_at = NOW() WHERE id = $3`,
		status, retryCount, id,
	)
	return err
}

// fetches a webhook by id
func (r *WebhookRepo) GetByID(ctx context.Context, id uuid.UUID) (*Webhook, error) {
	w := &Webhook{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, source, payload, target_url, status, retry_count, created_at, updated_at
		 FROM webhooks WHERE id = $1`, id,
	).Scan(&w.ID, &w.Source, &w.Payload, &w.TargetURL, &w.Status, &w.RetryCount, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return w, nil
}

// returns recent webhooks with pagination
func (r *WebhookRepo) ListRecent(ctx context.Context, limit, offset int) ([]Webhook, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM webhooks`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, source, payload, target_url, status, retry_count, created_at, updated_at
		 FROM webhooks ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var webhooks []Webhook
	for rows.Next() {
		var w Webhook
		if err := rows.Scan(&w.ID, &w.Source, &w.Payload, &w.TargetURL, &w.Status, &w.RetryCount, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, 0, err
		}
		webhooks = append(webhooks, w)
	}

	return webhooks, total, rows.Err()
}
