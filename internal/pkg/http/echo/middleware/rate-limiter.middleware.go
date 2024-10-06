package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// ApplyRateLimiter applies rate limiting middleware.
func ApplyRateLimiter(e *echo.Echo, rate rate.Limit) {
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(rate),
	}))
}
