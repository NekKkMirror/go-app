package middleware

import (
	"fmt"
	"time"

	spanutils "github.com/NekKkMirror/go-app/internal/pkg/utils/otel/span"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// EchoTracerMiddleware is an OpenTelemetry middleware for the Echo web framework.
func EchoTracerMiddleware(serviceName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			ctx := request.Context()

			ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(request.Header))

			spanName := c.Path()
			if spanName == "" {
				spanName = fmt.Sprintf("HTTP %s route not found", request.Method)
			}

			// Start the span with comprehensive options
			opts := []oteltrace.SpanStartOption{
				oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request)...),
				oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request)...),
				oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(serviceName, c.Path(), request)...),
				oteltrace.WithAttributes(
					attribute.String("http.method", request.Method),
					attribute.String("http.url", request.URL.String()),
					attribute.String("http.target", request.URL.Path),
					attribute.String("http.host", request.Host),
					attribute.String("http.scheme", request.URL.Scheme),
					attribute.String("http.server_name", serviceName),
					attribute.String("http.client_ip", c.RealIP()),
				),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
				oteltrace.WithNewRoot(),
			}
			ctx, span := otel.Tracer("echo-http").Start(ctx, spanName, opts...)
			defer span.End()

			c.SetRequest(request.WithContext(ctx))

			// Record start time to measure the request duration
			startTime := time.Now()

			// Call the next handler in the chain
			err := next(c)

			duration := time.Since(startTime)

			response := c.Response()

			spanutils.SetSpanAttributes(span, map[string]interface{}{
				"http.status_code":   response.Status,
				"http.user_agent":    request.UserAgent(),
				"http.response.size": response.Size,
				"http.duration":      duration.String(),
			})

			// Handle any errors that occurred during the request
			if err != nil {
				c.Error(err)

				var echoError *echo.HTTPError
				if errors.As(err, &echoError) {
					span.SetStatus(codes.Error, echoError.Message.(string))
					span.SetAttributes(
						attribute.String("echo.error", echoError.Message.(string)),
					)
				} else {
					span.SetStatus(codes.Error, err.Error())
					span.SetAttributes(attribute.String("echo.error", err.Error()))
				}
			} else {
				span.SetStatus(codes.Ok, "OK")
			}

			return err
		}
	}
}
