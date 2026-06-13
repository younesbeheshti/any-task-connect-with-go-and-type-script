package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
)

// Publisher publishes domain events to RabbitMQ.
type Publisher struct {
	client *Client
}

// NewPublisher creates a Publisher.
func NewPublisher(client *Client) *Publisher {
	return &Publisher{client: client}
}

// Publish sends an event to the configured exchange.
func (p *Publisher) Publish(ctx context.Context, routingKey string, payload any) error {
	ch, err := p.client.Channel()
	if err != nil {
		return err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	event := common.Event{
		ID:        uuid.NewString(),
		Name:      routingKey,
		Timestamp: time.Now().UTC(),
		Payload:   body,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return ch.PublishWithContext(ctx,
		p.client.Exchange(),
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now().UTC(),
			Body:         data,
		},
	)
}
