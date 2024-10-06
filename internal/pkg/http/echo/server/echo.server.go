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

// EchoConfig holds the configuration for the Echo server.
type EchoConfig struct {
	Host               string   `mapstructure:"host"`
	Port               string   `mapstructure:"port" validate:"required"`
	Development        string   `mapstructure:"development"`
	BasePath           string   `mapstructure:"basePath" validate:"required"`
	DebugErrorResponse bool     `mapstructure:"debugErrorResponse"`
	IgnoreLogUrls      []string `mapstructure:"ignoreLogUrls"`
	Timeout            string   `mapstructure:"timeout"`
}

// NewEchoServer creates and returns a new Echo instance.
func NewEchoServer() *echo.Echo {
	return echo.New()
}

// RunHttpServer runs the HTTP server and handles graceful shutdown on context cancellation.
func RunHttpServer(ctx context.Context, e *echo.Echo, log logger.ILogger, cfg *EchoConfig) error {
	configureServer(e, cfg)

	go func() {
		<-ctx.Done()
		shutdownServer(e, log, cfg, ctx)
	}()

	err := e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// configureServer sets the various server options such as timeouts and header sizes.
func configureServer(e *echo.Echo, cfg *EchoConfig) {
	e.Server.MaxHeaderBytes = MaxHeaderBytes
	e.Server.ReadTimeout = ReadTimeout
	e.Server.WriteTimeout = WriteTimeout
}

// shutdownServer gracefully shuts down the Echo server.
func shutdownServer(e *echo.Echo, log logger.ILogger, cfg *EchoConfig, ctx context.Context) {
	log.Infof("shutting down HTTP server on port: %s", cfg.Port)
	if err := e.Shutdown(ctx); err != nil {
		log.Errorf("(shutdown) error: %v", err)
		return
	}
	log.Info("server exited properly")
}

// ApplyVersioningFromHeader applies versioning to the Echo instance based on the "version" header.
func ApplyVersioningFromHeader(e *echo.Echo) {
	e.Pre(apiVersion)
}

// apiVersion is a middleware function that prefixes the request path with the version from the "version" header.
func apiVersion(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		version := req.Header.Get("version")
		req.URL.Path = fmt.Sprintf("/%s%s", version, req.URL.Path)
		return next(c)
	}
}

// RegisterGroupFunc registers a route group with the given name and builder function.
func RegisterGroupFunc(groupName string, e *echo.Echo, builder func(g *echo.Group)) *echo.Echo {
	builder(e.Group(groupName))
	return e
}
