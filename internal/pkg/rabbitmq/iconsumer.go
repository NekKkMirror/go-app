package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

//go:generate mockery --name IConsumer
type IConsumer[T any] interface {
	// Consume sets up the RabbitMQ consumer, processes incoming messages, and handles them using the provided dependencies.
	Consume(msg interface{}, dependencies T) error

	// openChannel opens a new RabbitMQ channel for consuming messages.
	openChannel(ctx context.Context, span trace.Span) (*amqp.Channel, error)

	// closeChannel closes the provided RabbitMQ channel and handles any errors that occur.
	closeChannel(channel *amqp.Channel, span trace.Span)

	// setupExchangeQueue sets up the RabbitMQ exchange, queue, and binds them together.
	setupExchangeQueue(channel *amqp.Channel, snakeTypeName string, span trace.Span) error

	// declareExchange declares an exchange in RabbitMQ with the given snakeTypeName and configuration.
	declareExchange(channel *amqp.Channel, snakeTypeName string, span trace.Span) error

	// declareQueue declares a queue in RabbitMQ with the given snakeTypeName as both the queue name and routing key.
	declareQueue(channel *amqp.Channel, snakeTypeName string, span trace.Span) (amqp.Queue, error)

	// bindQueue binds the declared queue to the specified exchange.
	bindQueue(channel *amqp.Channel, snakeTypeName string, span trace.Span) error

	// consumeMessages starts consuming messages from the specified queue.
	consumeMessages(channel *amqp.Channel, queueName string, span trace.Span) (<-chan amqp.Delivery, error)

	// handleMessage processes a single message and handles any errors that occur.
	handleMessage(ctx context.Context, span trace.Span, d amqp.Delivery, dependencies T) error

	// handleError logs the error, updates counters, and sets span attributes.
	handleError(logMessage string, err error, span trace.Span) error

	// IsConsumed checks if the message is consumed.
	IsConsumed(msg string) bool

	// Close closes the consumer and cleans up resources.
	Close() error

	// HealthCheck checks the health of the consumer.
	HealthCheck() bool

	// GetMetrics returns the metrics of consumed messages, error count and success count.
	GetMetrics() (int, int, int)
}
