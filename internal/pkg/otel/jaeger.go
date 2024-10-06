package otel

import (
	"context"
	"fmt"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/constants/environment"
	"github.com/NekKkMirror/go-app/internal/pkg/logger"
	"github.com/NekKkMirror/go-app/internal/pkg/utils/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp" // or otlptracegrpc
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

type OTLPConfig struct {
	Server      string `mapstructure:"server"`
	ServiceName string `mapstructure:"serviceName"`
	TracerName  string `mapstructure:"tracerName"`
}

// TracerProvider initializes and configures an OpenTelemetry TracerProvider with an OTLP exporter.
func TracerProvider(ctx context.Context, cfg *OTLPConfig, log logger.ILogger) (trace.Tracer, error) {
	if cfg.Server == "" {
		err := fmt.Errorf("failed to create exporter: signal: no such file or directory")
		log.Error("Failed to initialize TracerProvider", "error", err)
		return otel.Tracer(cfg.TracerName), err
	}

	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(cfg.Server))
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	env := common.GetEnv("APP_ENV", environment.Development)
	if env != environment.Production {
		env = environment.Development
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.ServiceName),
		attribute.String("environment", env),
		attribute.Float64("version", 1.0),
		attribute.String("team", "backend"),
		attribute.String("host", common.GetEnv("SERVICE_HOSTNAME", "localhost")),
	)

	// Create a TracerProvider with the configured exporter and resource
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter, tracesdk.WithBatchTimeout(5*time.Second), tracesdk.WithMaxExportBatchSize(100)),
		tracesdk.WithResource(resource),
	)

	go func() {
		<-ctx.Done()
		if err := tp.Shutdown(ctx); err != nil {
			log.Error("Failed to shutdown TracerProvider", "error", err)
		} else {
			log.Info("open-telemetry exited properly")
		}
	}()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	tracer := tp.Tracer(cfg.TracerName)

	_, span := tracer.Start(ctx, "TracerProvider Initialized",
		trace.WithAttributes(attribute.String("initialization", "success")),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	span.AddEvent("Tracer successfully initialized")

	log.Info("TracerProvider initialized")

	return tracer, nil
}
