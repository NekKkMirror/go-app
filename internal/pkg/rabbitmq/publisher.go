package rabbitmq

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/logger"
	"github.com/NekKkMirror/go-app/internal/pkg/otel"
	"github.com/cenkalti/backoff/v4"
	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Publisher struct {
	cfg               *RMQConfig
	ctx               context.Context
	conn              *amqp.Connection
	log               logger.ILogger
	tracer            trace.Tracer
	publishedMessages map[string]amqp.Publishing
	messageCounter    int
	errorCounter      int
	successCounter    int
}

// NewPublisher creates a new RabbitMQ publisher instance.
func NewPublisher(cfg *RMQConfig, ctx context.Context, conn *amqp.Connection, log logger.ILogger, tracer trace.Tracer) *Publisher {
	if err := ValidateConfig(cfg); err != nil {
		log.Error("Invalid configuration: ", err)
		return nil
	}
	return &Publisher{
		cfg:               cfg,
		ctx:               ctx,
		conn:              conn,
		log:               log,
		tracer:            tracer,
		publishedMessages: make(map[string]amqp.Publishing),
	}
}

// Publish marshals a message and publishes it to RabbitMQ with retry logic
func (p *Publisher) Publish(msg interface{}) error {
	data, err := jsoniter.Marshal(msg)
	if err != nil {
		return p.handleError("Error marshalling message", err, nil)
	}

	typeName := reflect.TypeOf(msg).Elem().Name()
	exchange := strcase.ToSnake(typeName)

	ctx, span := p.tracer.Start(p.ctx, typeName)
	defer span.End()

	headers := otel.InjectAMQPHeaders(ctx)
	headers["authenticated"] = true

	ch, err := p.openChannel(span)
	if err != nil {
		return err
	}
	defer p.closeChannel(ch, span)

	if err = p.declareExchange(ch, exchange, span); err != nil {
		return err
	}

	correlationID := p.getCorrelationID(ctx)
	msgID := uuid.NewV4().String()
	publishingMsg := p.createPublishingMessage(data, correlationID, msgID, exchange, headers, span)

	if err = p.publishMessageWithRetry(ch, exchange, publishingMsg, span); err != nil {
		return err
	}

	p.publishedMessages[msgID] = publishingMsg
	p.messageCounter++
	p.successCounter++

	p.log.Info("Message successfully published. Message ID: ", msgID)
	return nil
}

// openChannel opens a channel to RabbitMQ
func (p *Publisher) openChannel(span trace.Span) (*amqp.Channel, error) {
	channel, err := p.conn.Channel()
	if err != nil {
		return nil, p.handleError("Error opening channel", err, span)
	}
	return channel, nil
}

// closeChannel closes the channel to RabbitMQ
func (p *Publisher) closeChannel(channel *amqp.Channel, span trace.Span) {
	if err := channel.Close(); err != nil {
		_ = p.handleError("Error closing channel", err, span)
	}
}

// declareExchange declares an exchange in RabbitMQ
func (p *Publisher) declareExchange(channel *amqp.Channel, exchange string, span trace.Span) error {
	err := channel.ExchangeDeclare(
		exchange,   // name
		p.cfg.Kind, // kind
		true,       // durable
		false,      // auto-delete
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return p.handleError("Error declaring exchange", err, span)
	}
	return nil
}

// getCorrelationID retrieves the correlation ID from the context
func (p *Publisher) getCorrelationID(ctx context.Context) string {
	correlationID := ""
	if val := ctx.Value(echo.HeaderXCorrelationID); val != nil {
		if id, ok := val.(string); ok {
			correlationID = id
		}
	}
	if correlationID == "" {
		correlationID = uuid.NewV4().String()
	}
	return correlationID
}

// createPublishingMessage creates the message with the appropriate attributes.
func (p *Publisher) createPublishingMessage(data []byte, correlationID, msgID string, exchange string, headers otel.AMQPHeadersCarrier, span trace.Span) amqp.Publishing {
	msg := amqp.Publishing{
		Body:          data,
		ContentType:   p.cfg.ContentType,
		DeliveryMode:  amqp.Persistent,
		MessageId:     msgID,
		Timestamp:     time.Now(),
		CorrelationId: correlationID,
		Headers:       amqp.Table(headers),
		ReplyTo:       p.cfg.ReplyToQueue,
		Priority:      p.cfg.DefaultPriority,
		Expiration:    p.cfg.DefaultExpiration,
	}

	SetFixedPublishingSpanAttributes(span, msg, exchange, p.cfg.Kind, headers)
	return msg
}

// publishMessageWithRetry publishes the message with retry logic.
func (p *Publisher) publishMessageWithRetry(channel *amqp.Channel, exchange string, msg amqp.Publishing, span trace.Span) error {
	return backoff.Retry(func() error {
		err := channel.Publish(
			exchange, // exchange
			exchange, // routing key
			false,    // mandatory
			false,    // immediate
			msg,
		)
		if err != nil {
			p.log.Error("Error publishing message: ", err)
			span.SetAttributes(attribute.String("status", "error"), attribute.String("error.message", err.Error()))
			return err
		}
		return nil
	}, backoff.WithContext(backoff.NewExponentialBackOff(), p.ctx))
}

// handleError logs the error, updates counters, and sets span attributes
func (p *Publisher) handleError(logMessage string, err error, span trace.Span) error {
	p.log.Error(logMessage, err)
	p.errorCounter++
	if span != nil {
		span.SetAttributes(attribute.String("status", "error"), attribute.String("error.message", err.Error()))
	}
	return fmt.Errorf("%s: %w", logMessage, err)
}

// IsPublished checks if a message with the given ID has been published
func (p *Publisher) IsPublished(msgID string) bool {
	_, exists := p.publishedMessages[msgID]
	return exists
}

// Close closes the RabbitMQ connection
func (p *Publisher) Close() error {
	err := p.conn.Close()
	if err != nil {
		p.log.Error("Error closing connection: ", err)
		return fmt.Errorf("error closing connection: %w", err)
	}
	return nil
}

// HealthCheck - Checks the health of the publisher connection
func (p *Publisher) HealthCheck() bool {
	return p.conn.IsClosed() == false
}

// GetMetrics returns the total, successful, and failed message count
func (p *Publisher) GetMetrics() (total int, success int, failure int) {
	return p.messageCounter, p.successCounter, p.errorCounter
}
