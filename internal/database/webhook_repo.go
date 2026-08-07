package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// a row in the users table
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

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

// a row in the delivery_logs table
type DeliveryLog struct {
	ID            uuid.UUID `json:"id"`
	WebhookID     uuid.UUID `json:"webhook_id"`
	AttemptNumber int       `json:"attempt_number"`
	StatusCode    int       `json:"status_code"`
	ResponseBody  string    `json:"response_body"`
	CreatedAt     time.Time `json:"created_at"`
}

// crud operations for webhooks
type WebhookRepo struct {
	pool *pgxpool.Pool
}

// creates a new repo
func NewWebhookRepo(pool *pgxpool.Pool) *WebhookRepo {
	return &WebhookRepo{pool: pool}
}

// User Operations

func (r *WebhookRepo) CreateUser(ctx context.Context, email, passwordHash string) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, created_at) VALUES ($1, $2, $3, NOW())`,
		id, email, passwordHash,
	)
	return id, err
}

func (r *WebhookRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &u, nil
}

// inserts a webhook and returns its id
func (r *WebhookRepo) Insert(ctx context.Context, source string, payload map[string]interface{}, targetURL string, endpointID uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO webhooks (id, source, payload, target_url, endpoint_id, status, retry_count, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, 'PENDING', 0, NOW(), NOW())`,
		id, source, payload, targetURL, endpointID,
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

// returns recent webhooks with pagination, status filter, and search
func (r *WebhookRepo) ListRecent(ctx context.Context, limit, offset int, status string, search string, userID uuid.UUID) ([]Webhook, int, error) {
	var total int
	
	countQuery := `
		SELECT COUNT(*) FROM webhooks w
		JOIN endpoints e ON w.endpoint_id = e.id
		WHERE e.user_id = $1
		AND ($2 = '' OR w.status = $2) 
		AND ($3 = '' OR w.source ILIKE '%' || $3 || '%' OR w.id::text ILIKE '%' || $3 || '%')`
	err := r.pool.QueryRow(ctx, countQuery, userID, status, search).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT w.id, w.source, w.payload, w.target_url, w.status, w.retry_count, w.created_at, w.updated_at
		FROM webhooks w
		JOIN endpoints e ON w.endpoint_id = e.id
		WHERE e.user_id = $1
		AND ($2 = '' OR w.status = $2) 
		AND ($3 = '' OR w.source ILIKE '%' || $3 || '%' OR w.id::text ILIKE '%' || $3 || '%')
		ORDER BY w.created_at DESC
		LIMIT $4 OFFSET $5
	`
	rows, err := r.pool.Query(ctx, query, userID, status, search, limit, offset)
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

// inserts a delivery log for a webhook attempt
func (r *WebhookRepo) InsertDeliveryLog(ctx context.Context, webhookID uuid.UUID, attempt, statusCode int, responseBody string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO delivery_logs (webhook_id, attempt_number, status_code, response_body, created_at)
		 VALUES ($1, $2, $3, $4, NOW())`,
		webhookID, attempt, statusCode, responseBody,
	)
	return err
}

// fetches all delivery logs for a specific webhook
func (r *WebhookRepo) GetDeliveryLogs(ctx context.Context, webhookID uuid.UUID) ([]DeliveryLog, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, webhook_id, attempt_number, status_code, response_body, created_at
		 FROM delivery_logs WHERE webhook_id = $1 ORDER BY attempt_number ASC`, webhookID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []DeliveryLog
	for rows.Next() {
		var l DeliveryLog
		if err := rows.Scan(&l.ID, &l.WebhookID, &l.AttemptNumber, &l.StatusCode, &l.ResponseBody, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	return logs, rows.Err()
}

// a row in the endpoints table
type Endpoint struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	TargetURL string    `json:"target_url"`
	SecretKey string    `json:"secret_key"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// creates a new endpoint
func (r *WebhookRepo) CreateEndpoint(ctx context.Context, name, targetURL, secretKey string, userID uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO endpoints (id, name, target_url, secret_key, user_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, NOW())`,
		id, name, targetURL, secretKey, userID,
	)
	return id, err
}

// fetches an endpoint by ID and userID
func (r *WebhookRepo) GetEndpoint(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Endpoint, error) {
	e := &Endpoint{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, target_url, secret_key, user_id, created_at FROM endpoints WHERE id = $1 AND user_id = $2`, id, userID,
	).Scan(&e.ID, &e.Name, &e.TargetURL, &e.SecretKey, &e.UserID, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// fetches an endpoint by ID only (for ingestion)
func (r *WebhookRepo) GetEndpointForIngest(ctx context.Context, id uuid.UUID) (*Endpoint, error) {
	e := &Endpoint{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, target_url, secret_key, user_id, created_at FROM endpoints WHERE id = $1`, id,
	).Scan(&e.ID, &e.Name, &e.TargetURL, &e.SecretKey, &e.UserID, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// lists all endpoints for a user
func (r *WebhookRepo) ListEndpoints(ctx context.Context, userID uuid.UUID) ([]Endpoint, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, target_url, secret_key, user_id, created_at FROM endpoints WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []Endpoint
	for rows.Next() {
		var e Endpoint
		if err := rows.Scan(&e.ID, &e.Name, &e.TargetURL, &e.SecretKey, &e.UserID, &e.CreatedAt); err != nil {
			return nil, err
		}
		endpoints = append(endpoints, e)
	}
	return endpoints, rows.Err()
}
