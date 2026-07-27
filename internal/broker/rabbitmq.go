package broker

import (
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeWebhook      = "webhook.exchange"
	ExchangeRetry        = "webhook.retry.exchange"
	QueueProcess         = "webhook.process.queue"
	QueueRetry           = "webhook.retry.queue"
	QueueDLQ             = "webhook.dlq.queue"
	RoutingKeyProcess    = "process"
	RoutingKeyRetry      = "retry"
	RoutingKeyDLQ        = "dlq"
)

// connects to rabbitmq with retries
func Connect(url string) (*amqp.Connection, error) {
	const maxRetries = 10
	const retryDelay = 3 * time.Second

	var conn *amqp.Connection
	var err error

	for i := range maxRetries {
		conn, err = amqp.Dial(url)
		if err == nil {
			slog.Info("connected to RabbitMQ")
			return conn, nil
		}
		slog.Warn("failed to connect to RabbitMQ, retrying...",
			"attempt", i+1,
			"error", err,
		)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", maxRetries, err)
}

// sets up all exchanges, queues, and bindings
func DeclareTopology(ch *amqp.Channel) error {
	// main webhook exchange
	if err := ch.ExchangeDeclare(ExchangeWebhook, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare exchange %s: %w", ExchangeWebhook, err)
	}

	// retry exchange
	if err := ch.ExchangeDeclare(ExchangeRetry, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare exchange %s: %w", ExchangeRetry, err)
	}

	// processing queue
	if _, err := ch.QueueDeclare(QueueProcess, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare queue %s: %w", QueueProcess, err)
	}

	// bind processing queue
	if err := ch.QueueBind(QueueProcess, RoutingKeyProcess, ExchangeWebhook, false, nil); err != nil {
		return fmt.Errorf("bind queue %s: %w", QueueProcess, err)
	}

	// retry queue with dlx back to main exchange (expired msgs get reprocessed)
	retryArgs := amqp.Table{
		"x-dead-letter-exchange":    ExchangeWebhook,
		"x-dead-letter-routing-key": RoutingKeyProcess,
	}
	if _, err := ch.QueueDeclare(QueueRetry, true, false, false, false, retryArgs); err != nil {
		return fmt.Errorf("declare queue %s: %w", QueueRetry, err)
	}

	// bind retry queue
	if err := ch.QueueBind(QueueRetry, RoutingKeyRetry, ExchangeRetry, false, nil); err != nil {
		return fmt.Errorf("bind queue %s: %w", QueueRetry, err)
	}

	// dead letter queue for exhausted retries
	if _, err := ch.QueueDeclare(QueueDLQ, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare queue %s: %w", QueueDLQ, err)
	}

	// bind dlq
	if err := ch.QueueBind(QueueDLQ, RoutingKeyDLQ, ExchangeWebhook, false, nil); err != nil {
		return fmt.Errorf("bind queue %s: %w", QueueDLQ, err)
	}

	slog.Info("RabbitMQ topology declared successfully")
	return nil
}
