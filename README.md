# Siphon Gateway

A resilient, fault-tolerant webhook ingestion and retry engine built in Go. Siphon Gateway catches incoming third-party webhooks, safely queues payloads using RabbitMQ, and forwards them with exponential backoff and Dead Letter Queue (DLQ) support to protect downstream APIs from dropped events and outages.

## Architecture

```
Webhook Provider ──▶ Ingestion API ──▶ RabbitMQ ──▶ Worker Pool ──▶ Target URL
                         │                              │
                         ▼                              ▼
                     PostgreSQL ◀────────────────── Status Updates
                         │
                         ▼
                    SSE Stream ──▶ Dashboard (SvelteKit)
```

## Quick Start

```bash
# Clone and start all services
docker compose up --build

# Send a test webhook
curl -X POST http://localhost:8080/api/v1/webhook \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Source: test" \
  -d '{"event": "payment.completed", "amount": 4999}'

# List recent webhooks
curl http://localhost:8080/api/v1/webhooks

# Stream live events
curl -N http://localhost:8080/api/v1/events

# Replay a failed webhook
curl -X POST http://localhost:8080/api/v1/webhook/{id}/replay
```

## Services

| Service    | Port  | Description                        |
|------------|-------|------------------------------------|
| Gateway    | 8080  | Go backend (API + Workers + SSE)   |
| PostgreSQL | 5432  | Webhook persistence                |
| RabbitMQ   | 5672  | Message broker                     |
| RabbitMQ UI| 15672 | Management console (guest/guest)   |

## Configuration

Copy `.env.example` to `.env` and configure:

| Variable            | Default                                      | Description                     |
|---------------------|----------------------------------------------|---------------------------------|
| `TARGET_URL`        | `http://host.docker.internal:9999/webhook`   | Downstream delivery URL         |
| `WORKER_CONCURRENCY`| `5`                                          | Number of concurrent workers    |
| `POSTGRES_USER`     | `siphon`                                     | PostgreSQL username             |
| `POSTGRES_PASSWORD` | `siphon_secret`                              | PostgreSQL password             |
| `POSTGRES_DB`       | `siphon_gateway`                             | PostgreSQL database name        |

## Retry Strategy

Failed deliveries are retried with exponential backoff:

| Retry | Delay |
|-------|-------|
| 1     | 2s    |
| 2     | 4s    |
| 3     | 8s    |
| 4     | 16s   |
| 5     | 32s   |

After 5 retries, the webhook is moved to the Dead Letter Queue (DLQ) with status `FAILED_DLQ`.

## Project Structure

```
├── cmd/gateway/          # Application entry point
├── internal/
│   ├── api/              # HTTP handlers (ingest, list, replay, SSE)
│   ├── broker/           # RabbitMQ connection, topology, publisher
│   ├── config/           # Environment-based configuration
│   ├── database/         # PostgreSQL pool + webhook repository
│   ├── sse/              # Server-Sent Events broadcast hub
│   └── worker/           # Message consumer + delivery + retry logic
├── migrations/           # SQL schema migrations
├── Dockerfile            # Multi-stage Go build
└── docker-compose.yml    # Full stack orchestration
```

## API Endpoints

| Method | Path                          | Description                    |
|--------|-------------------------------|--------------------------------|
| POST   | `/api/v1/webhook`             | Ingest a webhook payload       |
| GET    | `/api/v1/webhooks`            | List webhooks (paginated)      |
| GET    | `/api/v1/events`              | SSE real-time event stream     |
| POST   | `/api/v1/webhook/{id}/replay` | Replay a failed webhook        |
| GET    | `/health`                     | Health check                   |

## License

MIT
