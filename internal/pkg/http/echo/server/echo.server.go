package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const (
	MaxHeaderBytes = 1 << 20
	ReadTimeout    = 15 * time.Second
	WriteTimeout   = 15 * time.Second
)

type EchoConfig struct {
	Host               string   `mapstructure:"host"`
	Port               string   `mapstructure:"port"     validate:"required"`
	Development        string   `mapstructure:"development"`
	BasePath           string   `mapstructure:"basePath" validate:"required"`
	DebugErrorResponse bool     `mapstructure:"debugErrorResponse"`
	IgnoreLogUrls      []string `mapstructure:"ignoreLogUrls"`
	Timeout            string   `mapstructure:"timeout"`
}

// NewEchoServer creates and returns a new instance of Echo web framework.
// Echo is a fast and minimalist Go web framework that provides high
// performance, minimal memory footprint, and extensive middleware support.
//
// This function initializes a new Echo instance with default settings.
//
// Return value:
// - *echo.Echo: A pointer to the newly created Echo instance.
func NewEchoServer() *echo.Echo {
	e := echo.New()
	return e
}

// RunHttpServer starts a new HTTP server using the provided Echo instance and configuration.
// It listens for incoming requests on the specified host and port, and handles them using the provided Echo instance.
// The server is configured with default settings, such as maximum header bytes, read and write timeouts.
//
// The function also starts a separate goroutine to handle graceful shutdown. When the provided context is canceled,
// the server is shut down gracefully, and any ongoing requests are allowed to complete.
//
// If an error occurs while starting the server, it is logged, and the function returns the error.
// If the error is not of type http.ErrServerClosed, it is returned.
//
// Parameters:
// - ctx: A context that can be used to cancel the server gracefully.
// - echo: An instance of Echo, which provides the web framework functionality.
// - log: An instance of ILogger, which is used for logging.
// - cfg: A pointer to an EchoConfig struct, which contains the server configuration.
//
// Return value:
// - error: An error that occurred while starting the server, or nil if the server started successfully.
func RunHttpServer(ctx context.Context, echo *echo.Echo, log logger.ILogger, cfg *EchoConfig) error {
	echo.Server.MaxHeaderBytes = MaxHeaderBytes
	echo.Server.ReadTimeout = ReadTimeout
	echo.Server.WriteTimeout = WriteTimeout

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Infof("shutting down http port: {%s}", cfg.Port)
				err := echo.Shutdown(ctx)
				if err != nil {
					log.Errorf("(shutdown) err: {%v}", err)
					return
				}
				log.Info("server exited properly")
				return
			}
		}
	}()

	err := echo.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// ApplyVersioningFromHeader adds support for versioning by inspecting the "version" header in incoming HTTP requests.
// The middleware function modifies the request's URL path to include the version specified in the header.
// This allows for easy routing to different versions of API endpoints.
//
// Parameters:
//   - echo: A pointer to the Echo instance, which provides the web framework functionality.
//
// Return value:
//   - None
func ApplyVersioningFromHeader(echo *echo.Echo) {
	echo.Pre(apiVersion)
}

// apiVersion is a middleware function for Echo that adds support for versioning
// by inspecting the "version" header in incoming HTTP requests.
//
// The middleware function modifies the request's URL path to include the version
// specified in the header. This allows for easy routing to different versions of
// API endpoints.
//
// Parameters:
//   - next: The next handler in the middleware chain. This function will be called
//     after the versioning logic has been applied.
//
// Return value:
//   - echo.HandlerFunc: A function that takes an echo.Context as a parameter and
//     returns an error (if any). This function should be used as middleware in
//     an Echo group or route.
func apiVersion(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		headers := req.Header

		apiVersion := headers.Get("version")

		req.URL.Path = fmt.Sprintf("%s%s", apiVersion, req.URL.Path)

		return next(c)
	}
}

// RegisterGroupFunc is a function that registers a new Echo group with the provided group name and builder function.
// The function applies the builder function to the newly created group, allowing for the addition of routes and middleware.
//
// Parameters:
// - groupName (string): The name of the group to be registered. This name is used as a prefix for all routes within the group.
// - echo (*echo.Echo): A pointer to the Echo instance, which provides the web framework functionality.
// - builder (func(g *echo.Group)): A function that takes a pointer to an Echo group and adds routes and middleware to it.
//
// Return value:
// - *echo.Echo: A pointer to the Echo instance, which now contains the newly registered group.
func RegisterGroupFunc(groupName string, echo *echo.Echo, builder func(g *echo.Group)) *echo.Echo {
	builder(echo.Group(groupName))

	return echo
}
