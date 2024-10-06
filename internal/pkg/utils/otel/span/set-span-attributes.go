package span

import (
	"encoding/json"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// SetSpanAttributes sets multiple attributes on a span
func SetSpanAttributes(span trace.Span, attributes map[string]interface{}) {
	for key, value := range attributes {
		switch v := value.(type) {
		case string:
			span.SetAttributes(attribute.String(key, v))
		case int:
			span.SetAttributes(attribute.Int(key, v))
		case int64:
			span.SetAttributes(attribute.Int64(key, v))
		case float64:
			span.SetAttributes(attribute.Float64(key, v))
		case bool:
			span.SetAttributes(attribute.Bool(key, v))
		default:
			jsonValue, err := json.Marshal(v)
			if err != nil {
				span.SetAttributes(attribute.String(key, "could not serialize"))
			} else {
				span.SetAttributes(attribute.String(key, string(jsonValue)))
			}
		}
	}
}
