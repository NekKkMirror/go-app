package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ApplyCORS is a middleware function for Echo framework that applies Cross-Origin Resource Sharing (CORS)
func ApplyCORS(e *echo.Echo, allowOrigins []string, allowMethods []string) {
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}

	if len(allowMethods) == 0 {
		allowMethods = []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE}
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: allowMethods,
	}))
}
