package rabbitmq

import (
	"context"

	"github.com/NekKkMirror/go-app/internal/pkg/otel"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

//go:generate mockery --name IPublisher
type IPublisher interface {
	// Publish marshals a message and publishes it to RabbitMQ with retry logic.
	Publish(msg interface{}) error

	// openChannel opens a new RabbitMQ channel for publishing messages.
	openChannel(span trace.Span) (*amqp.Channel, error)

	// closeChannel closes the provided RabbitMQ channel and handles any errors that occur.
	closeChannel(channel *amqp.Channel, span trace.Span)

	// declareExchange declares an exchange in RabbitMQ with the given exchange name and configuration.
	declareExchange(channel *amqp.Channel, exchange string, span trace.Span) error

	// getCorrelationID retrieves the correlation ID from the context.
	getCorrelationID(ctx context.Context) string

	// createPublishingMessage creates the AMQP publishing message with the appropriate attributes.
	createPublishingMessage(data []byte, correlationID, msgID, exchange string, headers otel.AMQPHeadersCarrier, span trace.Span) amqp.Publishing

	// publishMessageWithRetry publishes the message with retry logic.
	publishMessageWithRetry(channel *amqp.Channel, exchange string, message amqp.Publishing, span trace.Span) error

	// handleError logs the error, updates counters, and sets span attributes.
	handleError(logMessage string, err error, span trace.Span) error

	// IsPublished checks if the message is published.
	IsPublished(msg string) bool

	// Close closes the publisher and cleans up resources.
	Close() error

	// HealthCheck checks the health of the publisher.
	HealthCheck() bool

	// GetMetrics returns the metrics of published messages, error count, and success count.
	GetMetrics() (int, int, int)
}
