package testutil

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
)

// MockPublisher records published events for assertions.
type MockPublisher struct {
	mu     sync.Mutex
	Events []PublishedEvent
}

// PublishedEvent captures a published message.
type PublishedEvent struct {
	RoutingKey string
	Payload    any
}

func (m *MockPublisher) Publish(_ context.Context, routingKey string, payload any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Events = append(m.Events, PublishedEvent{RoutingKey: routingKey, Payload: payload})
	return nil
}

// LastEvent returns the most recently published event.
func (m *MockPublisher) LastEvent() (PublishedEvent, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.Events) == 0 {
		return PublishedEvent{}, false
	}
	return m.Events[len(m.Events)-1], true
}

// MockConsumer is a test double for event consumption.
type MockConsumer struct {
	Handlers map[string]func(ctx context.Context, event common.Event) error
}

// NewMockConsumer creates a MockConsumer.
func NewMockConsumer() *MockConsumer {
	return &MockConsumer{Handlers: make(map[string]func(ctx context.Context, event common.Event) error)}
}

// Subscribe registers a handler for a routing key.
func (c *MockConsumer) Subscribe(_ context.Context, routingKey string, handler func(ctx context.Context, event common.Event) error) error {
	c.Handlers[routingKey] = handler
	return nil
}

// Dispatch invokes a registered handler.
func (c *MockConsumer) Dispatch(ctx context.Context, routingKey string, payload any) error {
	handler, ok := c.Handlers[routingKey]
	if !ok {
		return nil
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return handler(ctx, common.Event{Name: routingKey, Payload: body})
}
