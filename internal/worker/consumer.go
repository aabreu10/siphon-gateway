package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/aabreu10/siphon-gateway/internal/broker"
	"github.com/aabreu10/siphon-gateway/internal/database"
	"github.com/aabreu10/siphon-gateway/internal/sse"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	maxRetries      = 5
	deliveryTimeout = 10 * time.Second
	// Backwards-compatible alias (typo); prefer deliveryTimeout.
	deliveryTimout = deliveryTimeout
)

// launches n workers that consume and deliver webhooks
func StartWorkers(ctx context.Context, n int, conn *amqp.Connection, repo *database.WebhookRepo, pub *broker.Publisher, hub *sse.Hub) error {
	for i := range n {
		ch, err := conn.Channel()
		if err != nil {
			return fmt.Errorf("open channel for worker %d: %w", i, err)
		}

		// prefetch 1 msg at a time
		if err := ch.Qos(1, 0, false); err != nil {
			return fmt.Errorf("set QoS for worker %d: %w", i, err)
		}

		deliveries, err := ch.Consume(
			broker.QueueProcess, // queue
			fmt.Sprintf("worker-%d", i), // consumer tag
			false, // auto-ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			return fmt.Errorf("start consuming for worker %d: %w", i, err)
		}

		go runWorker(ctx, i, deliveries, repo, pub, hub)
	}

	slog.Info("worker pool started", "concurrency", n)
	return nil
}

func runWorker(ctx context.Context, id int, deliveries <-chan amqp.Delivery, repo *database.WebhookRepo, pub *broker.Publisher, hub *sse.Hub) {
	logger := slog.With("worker_id", id)
	logger.Info("worker started")

	for {
		select {
		case <-ctx.Done():
			logger.Info("worker shutting down")
			return
		case d, ok := <-deliveries:
			if !ok {
				logger.Info("delivery channel closed")
				return
			}
			processDelivery(ctx, logger, d, repo, pub, hub)
		}
	}
}

func processDelivery(ctx context.Context, logger *slog.Logger, d amqp.Delivery, repo *database.WebhookRepo, pub *broker.Publisher, hub *sse.Hub) {
	var msg broker.Message
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		logger.Error("failed to unmarshal message", "error", err)
		// reject, don't requeue
		_ = d.Nack(false, false)
		return
	}

	logger = logger.With(
		"webhook_id", msg.WebhookID,
		"retry_count", msg.RetryCount,
		"target_url", msg.TargetURL,
	)

	// attempt delivery
	statusCode, err := deliver(msg.Payload, msg.TargetURL)

	if err == nil && statusCode >= 200 && statusCode < 300 {
		// success
		logger.Info("delivery successful", "status_code", statusCode)
		if dbErr := repo.UpdateStatus(ctx, msg.WebhookID, "SUCCESS", msg.RetryCount); dbErr != nil {
			logger.Error("failed to update webhook status to SUCCESS", "error", dbErr)
		}
		hub.Broadcast(&sse.Event{
			Type: "success",
			Webhook: buildWebhookResponse(msg, "SUCCESS"),
		})
		_ = d.Ack(false)
		return
	}

	// failure
	if err != nil {
		logger.Warn("delivery failed with error", "error", err)
	} else {
		logger.Warn("delivery failed with non-2xx status", "status_code", statusCode)
	}

	msg.RetryCount++

	if msg.RetryCount <= maxRetries {
		// retry with backoff
		logger.Info("scheduling retry", "next_retry", msg.RetryCount)
		if dbErr := repo.UpdateStatus(ctx, msg.WebhookID, "PENDING", msg.RetryCount); dbErr != nil {
			logger.Error("failed to update retry count", "error", dbErr)
		}
		if pubErr := pub.PublishToRetry(ctx, &msg); pubErr != nil {
			logger.Error("failed to publish to retry queue", "error", pubErr)
			// requeue the current delivery so we don't drop the webhook if RabbitMQ is unavailable
			_ = d.Nack(false, true)
			return
		}
		hub.Broadcast(&sse.Event{
			Type: "retrying",
			Webhook: buildWebhookResponse(msg, "PENDING"),
		})
	} else {
		// retries exhausted, send to dlq
		logger.Warn("max retries exhausted, moving to DLQ")
		if dbErr := repo.UpdateStatus(ctx, msg.WebhookID, "FAILED_DLQ", msg.RetryCount); dbErr != nil {
			logger.Error("failed to update webhook status to FAILED_DLQ", "error", dbErr)
		}
		if pubErr := pub.PublishToDLQ(ctx, &msg); pubErr != nil {
			logger.Error("failed to publish to DLQ", "error", pubErr)
		}
		hub.Broadcast(&sse.Event{
			Type: "failed_dlq",
			Webhook: buildWebhookResponse(msg, "FAILED_DLQ"),
		})
	}

	_ = d.Ack(false)
}

// posts the payload to the target url, returns status code
func deliver(payload map[string]interface{}, targetURL string) (int, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("marshal payload: %w", err)
	}

	client := &http.Client{Timeout: deliveryTimout}
	resp, err := client.Post(targetURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return 0, fmt.Errorf("HTTP POST to %s: %w", targetURL, err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

// builds a webhook map for sse events
func buildWebhookResponse(msg broker.Message, status string) map[string]interface{} {
	return map[string]interface{}{
		"id":          msg.WebhookID,
		"source":      msg.Source,
		"payload":     msg.Payload,
		"target_url":  msg.TargetURL,
		"status":      status,
		"retry_count": msg.RetryCount,
	}
}
