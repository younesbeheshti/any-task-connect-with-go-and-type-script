package common

import (
	"context"
	"encoding/json"
	"time"
)

// Event names for RabbitMQ routing keys.
const (
	EventTaskCreated           = "task.created"
	EventTaskAssigned          = "task.assigned"
	EventTaskStarted           = "task.started"
	EventTaskCompleted         = "task.completed"
	EventTaskVerified          = "task.verified"
	EventPaymentReleased       = "payment.released"
	EventApplicationSubmitted  = "application.submitted"
	EventNotificationCreated   = "notification.created"
	EventReviewCreated         = "review.created"
)

// Event is the base envelope for domain events.
type Event struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// Publisher publishes domain events to the message broker.
type Publisher interface {
	Publish(ctx context.Context, routingKey string, payload any) error
}

// Consumer handles incoming domain events.
type Consumer interface {
	Subscribe(ctx context.Context, routingKey string, handler func(ctx context.Context, event Event) error) error
}
