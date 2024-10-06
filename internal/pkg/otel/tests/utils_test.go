package tests

import (
	"context"
	"net/http"
	"testing"

	myotel "github.com/NekKkMirror/go-app/internal/pkg/otel"
	"github.com/NekKkMirror/go-app/internal/pkg/otel/mocks"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
)

func TestAMQPHeadersCarrier(t *testing.T) {
	headers := make(myotel.AMQPHeadersCarrier)
	headers.Set("test-key", "test-value")

	assert.Equal(t, "test-value", headers.Get("test-key"))
	assert.Equal(t, "", headers.Get("nonexistent-key"))

	keys := headers.Keys()
	assert.Contains(t, keys, "test-key")
	assert.NotContains(t, keys, "nonexistent-key")
}

func TestInjectAMQPHeaders(t *testing.T) {
	ctx := context.Background()
	headers := myotel.InjectAMQPHeaders(ctx)

	assert.NotNil(t, headers)
}

func TestExtractAMQPHeaders(t *testing.T) {
	ctx := context.Background()
	headers := make(myotel.AMQPHeadersCarrier)
	headers.Set("test-key", "test-value")

	newCtx := myotel.ExtractAMQPHeaders(ctx, headers)
	assert.NotNil(t, newCtx)
}

func TestStartSpanFromContext(t *testing.T) {
	tracerProvider := mocks.NewMockTracerProvider()
	otel.SetTracerProvider(tracerProvider)

	ctx := context.Background()
	ctx, span := myotel.StartSpanFromContext(ctx, "test-span")

	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
}

func TestEndSpan(t *testing.T) {
	tracerProvider := mocks.NewMockTracerProvider()
	otel.SetTracerProvider(tracerProvider)

	ctx := context.Background()
	_, span := myotel.StartSpanFromContext(ctx, "test-span")

	myotel.EndSpan(span, nil)
}

func TestAddSpanAttributes(t *testing.T) {
	tracerProvider := mocks.NewMockTracerProvider()
	otel.SetTracerProvider(tracerProvider)

	ctx := context.Background()
	_, span := myotel.StartSpanFromContext(ctx, "test-span")

	attributes := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	myotel.AddSpanAttributes(span, attributes)
}

func TestTraceHTTPRequest(t *testing.T) {
	tracerProvider := mocks.NewMockTracerProvider()
	otel.SetTracerProvider(tracerProvider)

	ctx := context.Background()
	req, err := http.NewRequest("GET", "http://test.com", nil)
	assert.NoError(t, err)

	tracedReq := myotel.TraceHTTPRequest(ctx, req)

	assert.NotNil(t, tracedReq)
	assert.Equal(t, req, tracedReq)
}

func TestExtractAndStartSpanFromHTTPRequest(t *testing.T) {
	tracerProvider := mocks.NewMockTracerProvider()
	otel.SetTracerProvider(tracerProvider)

	ctx := context.Background()
	req, err := http.NewRequest("GET", "http://test.com", nil)
	assert.NoError(t, err)

	newCtx, span := myotel.ExtractAndStartSpanFromHTTPRequest(ctx, req, "test-span")
	assert.NotNil(t, newCtx)
	assert.NotNil(t, span)
}

func TestAddEventToSpan(t *testing.T) {
	tracerProvider := mocks.NewMockTracerProvider()
	otel.SetTracerProvider(tracerProvider)

	ctx := context.Background()
	_, span := myotel.StartSpanFromContext(ctx, "test-span")

	attributes := map[string]interface{}{
		"event-attr1": "value1",
		"event-attr2": "value2",
	}
	myotel.AddEventToSpan(span, "test-event", attributes)
}
