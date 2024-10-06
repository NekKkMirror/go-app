package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// CorrelationIdMiddleware adds a correlation ID to each HTTP request.
// If a correlation ID is already present in the request header, it will be used.
// Otherwise, a new UUID will be generated. The correlation ID will be added
// to the response header and stored in the request context for further use.
func CorrelationIdMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		const headerXCorrelationID = echo.HeaderXCorrelationID

		req := c.Request()
		id := getCorrelationID(req.Header.Get(headerXCorrelationID))

		c.Response().Header().Set(headerXCorrelationID, id)
		c.SetRequest(req.WithContext(context.WithValue(req.Context(), headerXCorrelationID, id)))

		return next(c)
	}
}

// getCorrelationID returns the correlation ID from the header or generates a new one.
func getCorrelationID(headerID string) string {
	if headerID == "" {
		return uuid.New().String()
	}
	return headerID
}
