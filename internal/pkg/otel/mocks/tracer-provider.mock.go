package mocks

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type MockTracerProvider struct {
	trace.TracerProvider
}

func NewMockTracerProvider() *MockTracerProvider {
	return &MockTracerProvider{noop.NewTracerProvider()}
}

func (tp *MockTracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return &MockTracer{}
}

type MockTracer struct {
	trace.Tracer
}

func (t *MockTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return ctx, &MockSpan{}
}

type MockSpan struct {
	trace.Span
}

func (s *MockSpan) End(opts ...trace.SpanEndOption) {}

func (s *MockSpan) AddLink(link trace.Link) {}

func (s *MockSpan) SetName(name string) {}

func (s *MockSpan) SetAttributes(attrs ...attribute.KeyValue) {}

func (s *MockSpan) TracerProvider() trace.TracerProvider {
	return &MockTracerProvider{}
}

func (s *MockSpan) SpanContext() trace.SpanContext {
	return trace.NewSpanContext(trace.SpanContextConfig{})
}

func (s *MockSpan) IsRecording() bool {
	return false
}

func (s *MockSpan) RecordError(err error, opts ...trace.EventOption) {}

func (s *MockSpan) AddEvent(name string, opts ...trace.EventOption) {}

func (s *MockSpan) AddEventWithTimestamp(timestamp time.Time, name string, opts ...trace.EventOption) {
}

func (s *MockSpan) SetStatus(code codes.Code, description string) {}
