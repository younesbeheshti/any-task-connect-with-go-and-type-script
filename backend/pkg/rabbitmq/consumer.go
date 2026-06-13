package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"go.uber.org/zap"
)

// Consumer subscribes to RabbitMQ routing keys.
type Consumer struct {
	client *Client
	log    *zap.Logger
}

// NewConsumer creates a Consumer.
func NewConsumer(client *Client, log *zap.Logger) *Consumer {
	return &Consumer{client: client, log: log}
}

// Subscribe registers a handler for a routing key.
func (c *Consumer) Subscribe(ctx context.Context, routingKey string, handler func(ctx context.Context, event common.Event) error) error {
	ch, err := c.client.Channel()
	if err != nil {
		return err
	}

	queueName := queueNameForRoutingKey(routingKey)
	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	if err := ch.QueueBind(q.Name, routingKey, c.client.Exchange(), false, nil); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				c.handleDelivery(ctx, msg, handler)
			}
		}
	}()

	c.log.Info("consumer subscribed", zap.String("routing_key", routingKey))
	return nil
}

func (c *Consumer) handleDelivery(ctx context.Context, msg amqp.Delivery, handler func(ctx context.Context, event common.Event) error) {
	var event common.Event
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		c.log.Error("unmarshal event", zap.Error(err))
		_ = msg.Nack(false, false)
		return
	}

	if err := handler(ctx, event); err != nil {
		c.log.Error("handle event", zap.Error(err), zap.String("event", event.Name))
		_ = msg.Nack(false, true)
		return
	}

	_ = msg.Ack(false)
}
