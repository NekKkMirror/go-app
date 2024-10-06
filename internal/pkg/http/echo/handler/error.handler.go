package handler

import (
	"net/http"

	"github.com/NekKkMirror/go-app/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// SetErrorHandler sets the HTTP error handler.
func SetErrorHandler(e *echo.Echo, log logger.ILogger, debug bool) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		var he *echo.HTTPError
		if errors.As(err, &he) {
			code = he.Code
		}

		logError := errors.Wrap(err, "HTTP error")
		log.Errorf("HTTP error (%d): %s", code, logError)

		if debug {
			c.JSON(code, map[string]interface{}{
				"error":   http.StatusText(code),
				"details": err.Error(),
			})
		} else {
			c.JSON(code, map[string]interface{}{
				"error": http.StatusText(code),
			})
		}
	}
}
