# Siphon Gateway

A resilient, production-grade webhook ingestion and retry engine built in Go. Siphon Gateway acts as a secure buffer between external API providers (like Stripe or GitHub) and your main application. It catches incoming webhooks, instantly acknowledges them, and safely queues payloads using RabbitMQ. Background workers then forward the webhooks to your target APIs with exponential backoff and Dead Letter Queue (DLQ) support, protecting your application from dropped events during deployments or outages.

## Key Features

- 🚀 **High-Performance Ingestion:** Built in Go using `chi` for minimal latency and high concurrency.
- 🔄 **Guaranteed Delivery:** RabbitMQ-backed asynchronous delivery with exponential backoff and DLQ routing.
- 🛡️ **Production-Ready Security:** 
  - Cross-domain `HttpOnly` JWT cookie authentication.
  - Strict CSRF mitigation (Content-Type enforcement).
  - HSTS and IP-based Rate Limiting on authentication routes.
  - Strict payload size limits (5MB) to prevent memory exhaustion.
- 🎛️ **Dynamic Endpoints:** Configure multiple target URLs dynamically via the database instead of hardcoding a single destination.
- 📊 **Real-Time Dashboard:** A responsive, modern dashboard built with **SvelteKit 5** and TailwindCSS, featuring live delivery metrics via Server-Sent Events (SSE).
- 🧪 **Built-in Simulator:** Test your webhook flows directly from the dashboard.

## Architecture

```text
Webhook Provider ──▶ Ingestion API ──▶ RabbitMQ ──▶ Worker Pool ──▶ Target URL
                         │                              │
                         ▼                              ▼
                     PostgreSQL ◀────────────────── Status Updates
                         │
                         ▼
                    SSE Stream ──▶ Dashboard (SvelteKit)
```

## Tech Stack

- **Backend API & Workers:** Go (Golang)
- **Frontend Dashboard:** SvelteKit 5, TypeScript, TailwindCSS
- **Message Broker:** RabbitMQ
- **Persistent Storage:** PostgreSQL

## API Endpoints

### Authentication
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/signup` | Create an admin account (Sets HttpOnly Cookie) |
| POST | `/api/v1/auth/login` | Authenticate and set HttpOnly Cookie |
| POST | `/api/v1/auth/logout` | Clear authentication cookie |

### Webhooks & Ingestion
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/ingest` | Ingest webhook (Requires `INGEST_API_KEY`) |
| GET | `/api/v1/webhooks` | List recent webhooks (Protected via Cookie) |
| GET | `/api/v1/webhook/{id}/logs` | Fetch delivery attempts for a webhook |
| POST | `/api/v1/webhook/{id}/replay` | Manually replay a failed DLQ webhook |

### Real-Time & Endpoints
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/events` | SSE real-time stream (Protected via Cookie) |
| GET | `/api/v1/endpoints` | List registered target URLs |
| POST | `/api/v1/endpoints` | Register a new target URL |

## Retry Strategy

Failed deliveries are retried automatically by the background workers using exponential backoff:

| Retry Attempt | Delay |
|---------------|-------|
| 1 | 2s |
| 2 | 4s |
| 3 | 8s |
| 4 | 16s |
| 5 | 32s |

After 5 failed retries, the webhook is moved to the Dead Letter Queue (DLQ) and marked as `FAILED_DLQ` in the database, where it can be manually replayed later via the dashboard.

## License

MIT
