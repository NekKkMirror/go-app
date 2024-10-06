package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/NekKkMirror/go-app/internal/pkg/otel"
	spanutils "github.com/NekKkMirror/go-app/internal/pkg/utils/otel/span"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

// ValidateConfig validates the configuration shared between Publisher and Consumer
func ValidateConfig(cfg *RMQConfig) error {
	if cfg.Kind == "" {
		return fmt.Errorf("invalid configuration: Kind, Exchange, and RoutingKey must be provided")
	}
	if cfg.MaxRetries < 0 {
		return fmt.Errorf("invalid configuration: MaxRetries cannot be negative")
	}
	return nil
}

// SetFixedPublishingSpanAttributes sets predefined attributes on a span for Publisher
func SetFixedPublishingSpanAttributes(span trace.Span, publishingMsg amqp.Publishing, exchange string, kind string, headers otel.AMQPHeadersCarrier) {
	headersJSON, err := json.Marshal(headers)
	if err != nil {
		headersJSON = []byte("could not serialize headers")
	}

	spanutils.SetSpanAttributes(span, map[string]interface{}{
		"status":           "pending",
		"message-id":       publishingMsg.MessageId,
		"correlation-id":   publishingMsg.CorrelationId,
		"exchange":         exchange,
		"kind":             kind,
		"content-type":     publishingMsg.ContentType,
		"timestamp":        publishingMsg.Timestamp.String(),
		"body":             string(publishingMsg.Body),
		"headers":          string(headersJSON),
		"reply-to":         publishingMsg.ReplyTo,
		"priority":         publishingMsg.Priority,
		"expiration":       publishingMsg.Expiration,
		"app-id":           publishingMsg.AppId,
		"user-id":          publishingMsg.UserId,
		"type":             publishingMsg.Type,
		"delivery-mode":    publishingMsg.DeliveryMode,
		"content-encoding": publishingMsg.ContentEncoding,
	})
}

// SetFixedConsumingSpanAttributes sets predefined attributes on a span for Consumer
func SetFixedConsumingSpanAttributes(span trace.Span, deliveryMsg amqp.Delivery, queueName string) {
	headersJSON, err := json.Marshal(deliveryMsg.Headers)
	if err != nil {
		headersJSON = []byte("could not serialize headers")
	}

	spanutils.SetSpanAttributes(span, map[string]interface{}{
		"status":         "pending",
		"message-id":     deliveryMsg.MessageId,
		"correlation-id": deliveryMsg.CorrelationId,
		"exchange":       deliveryMsg.Exchange,
		"queue":          queueName,
		"content-type":   deliveryMsg.ContentType,
		"timestamp":      deliveryMsg.Timestamp.String(),
		"body":           string(deliveryMsg.Body),
		"headers":        string(headersJSON),
		"reply-to":       deliveryMsg.ReplyTo,
		"priority":       deliveryMsg.Priority,
		"delivery-tag":   deliveryMsg.DeliveryTag,
		"redelivered":    deliveryMsg.Redelivered,
		"consumer-tag":   deliveryMsg.ConsumerTag,
	})
}
