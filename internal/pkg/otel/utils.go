package otel

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// AMQPHeadersCarrier is a custom carrier for AMQP headers to enable context propagation.
type AMQPHeadersCarrier map[string]interface{}

// Get retrieves the value associated with the given key from the AMQPHeadersCarrier.
func (a AMQPHeadersCarrier) Get(key string) string {
	v, ok := a[key]
	if !ok {
		return ""
	}
	return v.(string)
}

// Set associates the given key and value in the AMQPHeadersCarrier.
func (a AMQPHeadersCarrier) Set(key, value string) {
	a[key] = value
}

// Keys returns a slice of keys from the AMQPHeadersCarrier.
func (a AMQPHeadersCarrier) Keys() []string {
	keys := make([]string, 0, len(a))
	for k := range a {
		keys = append(keys, k)
	}
	return keys
}

// InjectAMQPHeaders injects the current context's trace and span information into an AMQPHeadersCarrier.
func InjectAMQPHeaders(ctx context.Context) AMQPHeadersCarrier {
	h := make(AMQPHeadersCarrier)
	otel.GetTextMapPropagator().Inject(ctx, h)
	return h
}

// ExtractAMQPHeaders extracts trace and span information from the provided AMQP headers and returns a new context with the extracted information.
func ExtractAMQPHeaders(ctx context.Context, headers AMQPHeadersCarrier) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, headers)
}

// EndSpan ends the provided span, capturing any error state.
func EndSpan(span oteltrace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Ok, "OK")
	}
	span.End()
}

// AddSpanAttributes adds attributes to the provided span.
func AddSpanAttributes(span oteltrace.Span, attributes map[string]interface{}) {
	for k, v := range attributes {
		span.SetAttributes(attribute.String(k, v.(string)))
	}
}

// TraceHTTPRequest injects tracing information into an HTTP request.
func TraceHTTPRequest(ctx context.Context, req *http.Request) *http.Request {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	return req.WithContext(ctx)
}

// ExtractAndStartSpanFromHTTPRequest extracts tracing information from an HTTP request and starts a new span.
func ExtractAndStartSpanFromHTTPRequest(ctx context.Context, req *http.Request, spanName string) (context.Context, oteltrace.Span) {
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(req.Header))
	return StartSpanFromContext(ctx, spanName)
}

// StartSpanFromContext starts a new span from the provided context.
func StartSpanFromContext(ctx context.Context, spanName string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	tracer := otel.Tracer("otel-utils")
	ctx, span := tracer.Start(ctx, spanName, opts...)
	return ctx, span
}

// AddEventToSpan adds an event to the provided span.
func AddEventToSpan(span oteltrace.Span, eventName string, attributes map[string]interface{}) {
	attrs := make([]attribute.KeyValue, 0, len(attributes))
	for k, v := range attributes {
		attrs = append(attrs, attribute.String(k, v.(string)))
	}
	span.AddEvent(eventName, oteltrace.WithAttributes(attrs...))
}
