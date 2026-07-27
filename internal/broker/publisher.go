package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// payload for rabbitmq messages
type Message struct {
	WebhookID  uuid.UUID              `json:"webhook_id"`
	Payload    map[string]interface{} `json:"payload"`
	TargetURL  string                 `json:"target_url"`
	RetryCount int                    `json:"retry_count"`
	Source     string                 `json:"source"`
}

// wraps an amqp channel for publishing
type Publisher struct {
	ch *amqp.Channel
}

// creates a publisher
func NewPublisher(ch *amqp.Channel) *Publisher {
	return &Publisher{ch: ch}
}

// publishes to the processing queue
func (p *Publisher) PublishToProcess(ctx context.Context, msg *Message) error {
	return p.publish(ctx, ExchangeWebhook, RoutingKeyProcess, msg, "")
}

// publishes to retry queue with exponential backoff ttl
func (p *Publisher) PublishToRetry(ctx context.Context, msg *Message) error {
	delayMs := int(math.Pow(2, float64(msg.RetryCount))) * 1000
	expiration := fmt.Sprintf("%d", delayMs)
	return p.publish(ctx, ExchangeRetry, RoutingKeyRetry, msg, expiration)
}

// publishes to the dead letter queue
func (p *Publisher) PublishToDLQ(ctx context.Context, msg *Message) error {
	return p.publish(ctx, ExchangeWebhook, RoutingKeyDLQ, msg, "")
}

func (p *Publisher) publish(ctx context.Context, exchange, routingKey string, msg *Message, expiration string) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	pub := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         body,
	}

	if expiration != "" {
		pub.Expiration = expiration
	}

	return p.ch.PublishWithContext(ctx, exchange, routingKey, false, false, pub)
}
