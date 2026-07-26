# siphon-gateway
A resilient, fault-tolerant webhook ingestion and retry engine built in Go. Siphon Gateway catches incoming third-party webhooks, safely queues payloads using RabbitMQ, and forwards them with exponential backoff and Dead Letter Queue (DLQ) support to protect downstream APIs from dropped events and outages.
