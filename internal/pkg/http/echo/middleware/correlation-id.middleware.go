package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// CorrelationIdMiddleware is an Echo middleware function that adds a correlation ID to each HTTP request.
// If a correlation ID is already present in the request header, it will be used. Otherwise, a new UUID will be generated.
// The correlation ID will be added to the response header and also stored in the request context for further use.
func CorrelationIdMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		id := req.Header.Get(echo.HeaderXCorrelationID)
		if id == "" {
			id = uuid.NewV4().String()
		}

		c.Response().Header().Set(echo.HeaderXCorrelationID, id)
		newReq := req.WithContext(context.WithValue(req.Context(), echo.HeaderXCorrelationID, id))
		c.SetRequest(newReq)

		return next(c)
	}
}
