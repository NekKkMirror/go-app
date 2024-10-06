package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/NekKkMirror/go-app/internal/pkg/logger"
	"github.com/iancoleman/strcase"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Consumer[T any] struct {
	cfg              *RMQConfig
	ctx              context.Context
	conn             *amqp.Connection
	log              logger.ILogger
	tracer           trace.Tracer
	consumedMessages map[string]amqp.Delivery
	messageCounter   int
	errorCounter     int
	successCounter   int
	messageHandler   MessageHandler[T]
}

// NewConsumer creates a new RabbitMQ consumer instance with the provided configuration, context, connection, logger, tracer, and message handler.
func NewConsumer[T any](cfg *RMQConfig, ctx context.Context, conn *amqp.Connection, log logger.ILogger, tracer trace.Tracer, handler MessageHandler[T]) *Consumer[T] {
	if err := ValidateConfig(cfg); err != nil {
		log.Error("Invalid configuration: ", err)
		return nil
	}
	return &Consumer[T]{
		cfg:              cfg,
		ctx:              ctx,
		conn:             conn,
		log:              log,
		tracer:           tracer,
		consumedMessages: make(map[string]amqp.Delivery),
		messageHandler:   handler,
	}
}

// Consume sets up the RabbitMQ consumer, processes incoming messages, and handles them using the provided dependencies.
func (c *Consumer[T]) Consume(msg interface{}, dependencies T) error {
	typeName := reflect.TypeOf(msg).Name()
	snakeTypeName := strcase.ToSnake(typeName)

	ctx, span := c.tracer.Start(c.ctx, typeName)
	defer span.End()

	channel, err := c.openChannel(ctx, span)
	if err != nil {
		return err
	}
	defer c.closeChannel(channel, span)

	if err := c.setupExchangeQueue(channel, snakeTypeName, span); err != nil {
		return err
	}

	qname := fmt.Sprintf("%s_%s", snakeTypeName, "queue")

	msgs, err := c.consumeMessages(channel, qname, span)
	if err != nil {
		return err
	}

	for d := range msgs {
		if err := c.handleMessage(ctx, span, d, qname, dependencies); err != nil {
			return err
		}
	}

	return nil
}

// openChannel opens a new RabbitMQ channel for consuming messages.
func (c *Consumer[T]) openChannel(ctx context.Context, span trace.Span) (*amqp.Channel, error) {
	channel, err := c.conn.Channel()
	if err != nil {
		return nil, c.handleError("Error opening channel", err, span)
	}
	return channel, nil
}

// closeChannel closes the provided RabbitMQ channel and handles any errors that occur.
func (c *Consumer[T]) closeChannel(channel *amqp.Channel, span trace.Span) {
	if err := channel.Close(); err != nil {
		_ = c.handleError("Error closing channel", err, span)
	}
}

// setupExchangeQueue sets up the RabbitMQ exchange, queue, and binds them together.
func (c *Consumer[T]) setupExchangeQueue(channel *amqp.Channel, snakeTypeName string, span trace.Span) error {
	if err := c.declareExchange(channel, snakeTypeName, span); err != nil {
		return err
	}

	if _, err := c.declareQueue(channel, snakeTypeName, span); err != nil {
		return err
	}

	if err := c.bindQueue(channel, snakeTypeName, span); err != nil {
		return err
	}

	return nil
}

// declareExchange declares an exchange in RabbitMQ with the given snakeTypeName and configuration.
func (c *Consumer[T]) declareExchange(channel *amqp.Channel, snakeTypeName string, span trace.Span) error {
	err := channel.ExchangeDeclare(
		snakeTypeName, // name
		c.cfg.Kind,    // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return c.handleError("Error declaring exchange", err, span)
	}
	return nil
}

// declareQueue declares a queue in RabbitMQ with the given snakeTypeName as both the queue name and routing key.
func (c *Consumer[T]) declareQueue(channel *amqp.Channel, snakeTypeName string, span trace.Span) (amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		fmt.Sprintf("%s_%s", snakeTypeName, "queue"), // name
		true,  // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return queue, c.handleError("Error declaring queue", err, span)
	}
	return queue, nil
}

// bindQueue binds the queue to the exchange with the given snakeTypeName as both the queue name and routing key.
func (c *Consumer[T]) bindQueue(channel *amqp.Channel, snakeTypeName string, span trace.Span) error {
	err := channel.QueueBind(
		fmt.Sprintf("%s_%s", snakeTypeName, "queue"), // queue name
		snakeTypeName, // routing key
		snakeTypeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return c.handleError("Error binding queue to exchange", err, span)
	}
	return nil
}

// consumeMessages sets up a consumer for the specified queue and returns a channel of deliveries.
func (c *Consumer[T]) consumeMessages(channel *amqp.Channel, queueName string, span trace.Span) (<-chan amqp.Delivery, error) {
	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, c.handleError("Error consuming messages", err, span)
	}
	return msgs, nil
}

// handleMessage processes an incoming RabbitMQ message, unmarshals it, and handles it using the provided dependencies.
func (c *Consumer[T]) handleMessage(ctx context.Context, span trace.Span, d amqp.Delivery, qname string, dependencies T) error {
	var msg T
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		return c.handleError("Error unmarshalling message", err, span)
	}

	if err := c.messageHandler(ctx, msg, dependencies); err != nil {
		c.errorCounter++
		return c.handleError("Error handling message", err, span)
	}

	SetFixedConsumingSpanAttributes(span, d, qname)

	c.consumedMessages[d.MessageId] = d
	c.messageCounter++
	c.successCounter++
	span.SetAttributes(attribute.String("status", "consumed"))
	return nil
}

// handleError logs an error message, increments the error counter, and sets appropriate attributes on the provided tracing span.
func (c *Consumer[T]) handleError(logMessage string, err error, span trace.Span) error {
	c.log.Error(logMessage, err)
	c.errorCounter++
	if span != nil {
		span.SetAttributes(attribute.String("status", "error"), attribute.String("error.message", err.Error()))
	}
	return fmt.Errorf("%s: %w", logMessage, err)
}

// IsConsumed checks if a given message has been consumed by the RabbitMQ consumer.
func (c *Consumer[T]) IsConsumed(msg interface{}) bool {
	msgID := reflect.ValueOf(msg).FieldByName("MessageId").String()
	_, exists := c.consumedMessages[msgID]
	return exists
}

// Close closes the RabbitMQ connection and returns an error if any occur.
func (c *Consumer[T]) Close() error {
	err := c.conn.Close()
	if err != nil {
		c.log.Error("Error closing RabbitMQ connection: ", err)
		c.errorCounter++
		return fmt.Errorf("error closing RabbitMQ connection: %w", err)
	}
	c.log.Info("RabbitMQ connection closed successfully")
	return nil
}

// HealthCheck checks the health of the RabbitMQ connection by attempting to open a channel.
func (c *Consumer[T]) HealthCheck() bool {
	channel, err := c.conn.Channel()
	if err != nil {
		c.log.Error("Health check failed: ", err)
		return false
	}
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(channel)
	return true
}

// GetMetrics retrieves the current metrics of the RabbitMQ consumer.
func (c *Consumer[T]) GetMetrics() (int, int, int) {
	return c.messageCounter, c.errorCounter, c.successCounter
}
