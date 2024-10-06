package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/grafana"
	"github.com/NekKkMirror/go-app/internal/pkg/logger"
	"github.com/NekKkMirror/go-app/internal/pkg/otel"
	"github.com/NekKkMirror/go-app/internal/pkg/utils/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Metrics struct {
	HTTPRequestCount       *prometheus.CounterVec
	HTTPErrorCount         *prometheus.CounterVec
	HTTPRequestDuration    *prometheus.HistogramVec
	BusinessOperationCount *prometheus.CounterVec
	ActiveRequests         prometheus.Gauge
	CPUUsageMetric         prometheus.Gauge
	MemoryUsageMetric      prometheus.Gauge
	DiskUsageMetric        prometheus.Gauge
	Tracer                 trace.Tracer
	Logger                 logger.ILogger
}

func NewMetrics(ctx context.Context, logger logger.ILogger, jaegerConfig *otel.OTLPConfig) (*Metrics, error) {
	tracer, err := otel.TracerProvider(ctx, jaegerConfig, logger)
	if err != nil {
		return nil, err
	}

	requestCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_count",
			Help: "Total number of HTTP requests received.",
		},
		[]string{"method", "endpoint"},
	)

	errorCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_error_count",
			Help: "Number of HTTP errors by response status.",
		},
		[]string{"method", "endpoint", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response latency (seconds) for HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	businessOperationCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_operation_count",
			Help: "Count of business operations, labelled by operation type.",
		},
		[]string{"operation", "result"},
	)

	activeRequests := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "active_requests",
		Help: "Current number of active requests being processed.",
	})

	cpuUsageMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "Current CPU usage percentage.",
	})

	memoryUsageMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage_bytes",
		Help: "Current memory usage in bytes.",
	})

	diskUsageMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_usage_bytes",
		Help: "Current disk usage in bytes.",
	})

	prometheus.MustRegister(
		requestCount, errorCount, requestDuration,
		businessOperationCount, activeRequests,
		cpuUsageMetric, memoryUsageMetric, diskUsageMetric,
	)

	return &Metrics{
		HTTPRequestCount:       requestCount,
		HTTPErrorCount:         errorCount,
		HTTPRequestDuration:    requestDuration,
		BusinessOperationCount: businessOperationCount,
		ActiveRequests:         activeRequests,
		CPUUsageMetric:         cpuUsageMetric,
		MemoryUsageMetric:      memoryUsageMetric,
		DiskUsageMetric:        diskUsageMetric,
		Tracer:                 tracer,
		Logger:                 logger,
	}, nil
}

func (m *Metrics) ObserveHTTPRequestDuration(method, endpoint string, duration time.Duration) {
	m.HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	_, span := m.Tracer.Start(context.Background(), "ObserveHTTPRequestDuration")
	defer span.End()
	m.Logger.Infof("Duration for %s %s: %v seconds", method, endpoint, duration.Seconds())
}

func (m *Metrics) IncrementBusinessOperation(operation, result string) {
	m.BusinessOperationCount.WithLabelValues(operation, result).Inc()
	_, span := m.Tracer.Start(context.Background(), "IncrementBusinessOperation")
	defer span.End()
	m.Logger.Infof("Business operation %s resulted in %s", operation, result)
}

func StartMetricServer(defaultAddr string, grafanaConfig *grafana.Config) error {
	addr := common.GetEnv("METRICS_SERVER_ADDRESS", defaultAddr)
	http.Handle("/metrics", promhttp.Handler())
	logrus.Printf("Starting metrics server at %s", addr)

	prometheusURL := common.GetEnv("PROMETHEUS_URL", "http://localhost:9090")

	grafanaService := grafana.New(grafanaConfig)
	err := grafanaService.AddDataSource("Prometheus", "prometheus", prometheusURL)
	if err != nil {
		logrus.Errorf("Failed to add data source to Grafana: %v", err)
	}

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		return errors.New("Metrics server error: " + err.Error())
	}
	return nil
}
