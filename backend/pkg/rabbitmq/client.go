package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/younesbeheshti/any-task-connect/backend/configs"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"go.uber.org/zap"
)

const (
	exchangeType = "topic"
	reconnectDelay = 5 * time.Second
)

// Client manages RabbitMQ connection, exchange, and lifecycle.
type Client struct {
	cfg        configs.RabbitMQConfig
	log        *zap.Logger
	conn       *amqp.Connection
	channel    *amqp.Channel
	mu         sync.RWMutex
	closed     bool
	reconnecting bool
}

// Connect establishes a RabbitMQ client with exchange setup.
func Connect(ctx context.Context, cfg configs.RabbitMQConfig, log *zap.Logger) (*Client, error) {
	c := &Client{cfg: cfg, log: log}
	if err := c.connect(ctx); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) connect(ctx context.Context) error {
	conn, err := amqp.Dial(c.cfg.URL)
	if err != nil {
		return fmt.Errorf("dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("open channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		c.cfg.Exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return fmt.Errorf("declare exchange: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.channel = ch
	c.mu.Unlock()

	log := c.log
	log.Info("rabbitmq connected",
		zap.String("exchange", c.cfg.Exchange),
	)

	go c.watchConnection()

	return nil
}

func (c *Client) watchConnection() {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	if conn == nil {
		return
	}

	notifyClose := conn.NotifyClose(make(chan *amqp.Error, 1))
	err := <-notifyClose
	if err == nil {
		return
	}

	c.mu.RLock()
	closed := c.closed
	c.mu.RUnlock()
	if closed {
		return
	}

	c.log.Warn("rabbitmq connection lost, reconnecting", zap.Error(err))
	c.reconnect()
}

func (c *Client) reconnect() {
	c.mu.Lock()
	if c.reconnecting || c.closed {
		c.mu.Unlock()
		return
	}
	c.reconnecting = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.reconnecting = false
		c.mu.Unlock()
	}()

	for {
		c.mu.RLock()
		closed := c.closed
		c.mu.RUnlock()
		if closed {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := c.connect(ctx)
		cancel()
		if err == nil {
			c.log.Info("rabbitmq reconnected")
			return
		}

		c.log.Warn("rabbitmq reconnect failed", zap.Error(err))
		time.Sleep(reconnectDelay)
	}
}

// Channel returns the active AMQP channel.
func (c *Client) Channel() (*amqp.Channel, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.channel == nil {
		return nil, fmt.Errorf("rabbitmq channel not available")
	}
	return c.channel, nil
}

// Exchange returns the configured exchange name.
func (c *Client) Exchange() string {
	return c.cfg.Exchange
}

// HealthCheck verifies RabbitMQ connectivity.
func (c *Client) HealthCheck(ctx context.Context) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	if conn == nil || conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection closed")
	}
	// Probe liveness by opening a dedicated short-lived channel rather than reusing
	// the shared channel. A channel-level exception (e.g. a passive declare of a
	// missing queue) closes the channel it runs on, so probing the shared channel
	// would break event publishing. Opening a channel already requires a live
	// connection and a responsive broker, making it a sufficient check.
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq channel open failed: %w", err)
	}
	_ = ch.Close()
	return nil
}

// Close closes the RabbitMQ connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// SetupEventQueues declares queues for known domain events.
func (c *Client) SetupEventQueues(ctx context.Context) error {
	ch, err := c.Channel()
	if err != nil {
		return err
	}

	events := []string{
		common.EventTaskCreated,
		common.EventTaskAssigned,
		common.EventTaskCompleted,
		common.EventPaymentReleased,
		common.EventNotificationCreated,
	}

	for _, routingKey := range events {
		queueName := queueNameForRoutingKey(routingKey)
		if _, err := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
			return fmt.Errorf("declare queue %s: %w", queueName, err)
		}
		if err := ch.QueueBind(queueName, routingKey, c.cfg.Exchange, false, nil); err != nil {
			return fmt.Errorf("bind queue %s: %w", queueName, err)
		}
	}

	c.log.Info("rabbitmq event queues configured", zap.Int("count", len(events)))
	return nil
}

func queueNameForRoutingKey(key string) string {
	return "taskbridge." + key
}
