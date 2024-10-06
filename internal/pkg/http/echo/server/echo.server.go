package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/constants/environment"
	"github.com/NekKkMirror/go-app/internal/pkg/http/echo/handler"
	"github.com/NekKkMirror/go-app/internal/pkg/http/echo/middleware"
	"github.com/NekKkMirror/go-app/internal/pkg/logger"
	"github.com/NekKkMirror/go-app/internal/pkg/utils/common"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

const (
	MaxHeaderBytes = 1 << 20
)

// EchoServer wraps the Echo instance and related configuration.
type EchoServer struct {
	Echo *echo.Echo
	Cfg  *EchoConfig
	Log  logger.ILogger
}

// EchoConfig holds the configuration for the Echo server.
type EchoConfig struct {
	Host               string   `mapstructure:"host" validate:"required"`
	Port               string   `mapstructure:"port" validate:"required"`
	BasePath           string   `mapstructure:"basePath" validate:"required"`
	DebugErrorResponse bool     `mapstructure:"debugErrorResponse"`
	IgnoreLogUrls      []string `mapstructure:"ignoreLogUrls"`
	Timeout            time.Duration

	CORSAllowOrigins []string `mapstructure:"cors_allow_origins"`
	CORSAllowMethods []string `mapstructure:"cors_allow_methods"`
	RateLimit        float64  `mapstructure:"rate_limit"`
}

// NewEchoServer creates and returns a configured EchoServer instance.
func NewEchoServer(cfg *EchoConfig, log logger.ILogger) *EchoServer {
	e := echo.New()
	if common.GetEnv("APP_ENV", environment.Development) == environment.Development {
		e.Debug = true
	}
	e.HideBanner = true

	e.Pre(apiVersion)

	server := &EchoServer{
		Echo: e,
		Cfg:  cfg,
		Log:  log,
	}
	server.configureServer()
	server.registerRoutes()
	server.applyMiddleware()
	handler.SetErrorHandler(e, log, cfg.DebugErrorResponse)

	return server
}

// Run runs the HTTP server and handles graceful shutdown on context cancellation.
func (s *EchoServer) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.shutdownServer(ctx)
	}()

	err := s.Echo.Start(fmt.Sprintf("%s:%s", s.Cfg.Host, s.Cfg.Port))
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// configureServer sets various server options, such as timeouts and header sizes.
func (s *EchoServer) configureServer() {
	s.Echo.Server.MaxHeaderBytes = MaxHeaderBytes
	s.Echo.Server.ReadTimeout = s.Cfg.Timeout
	s.Echo.Server.WriteTimeout = s.Cfg.Timeout
}

// shutdownServer gracefully shuts down the Echo server.
func (s *EchoServer) shutdownServer(ctx context.Context) {
	s.Log.Infof("shutting down HTTP server on port: %s", s.Cfg.Port)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Echo.Shutdown(shutdownCtx); err != nil {
		s.Log.Errorf("Graceful shutdown failed: %v", err)
	}
}

// registerRoutes registers all necessary routes including health-check.
func (s *EchoServer) registerRoutes() {
	s.Echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "UP"})
	})
}

// applyMiddleware applies the necessary middleware to the Echo instance.
func (s *EchoServer) applyMiddleware() {
	middleware.ApplySecurityHeaders(s.Echo)
	middleware.ApplyCORS(s.Echo, s.Cfg.CORSAllowOrigins, s.Cfg.CORSAllowMethods)
	middleware.ApplyRateLimiter(s.Echo, rate.Limit(s.Cfg.RateLimit))
	s.Echo.Use(middleware.CorrelationIdMiddleware)
	s.Echo.Use(echomiddleware.Logger())
	s.Echo.Use(echomiddleware.Recover())
	s.Echo.Use(echomiddleware.GzipWithConfig(echomiddleware.GzipConfig{
		Level: 5,
	}))

	// TODO below
	// Трассировка и мониторинг

	//tracer := otel.Tracer("echo-server")
	//s.Echo.Use(http.NewServerMiddleware(tracer))

	// Кэширование истекших экземпляров (например, Redis)
	// s.Echo.Use(redisMiddleware())

	// Идемпотентность
	// s.Echo.Use(idempotencyMiddleware())

	// Контентная договоренность
	// s.Echo.Use(contentNegotiationMiddleware())
}

// AddRoute add route with base path and handler function.
func (s *EchoServer) AddRoute(method, path string, handlerFunc echo.HandlerFunc) {
	fullPath := s.Cfg.BasePath + path
	s.Echo.Add(method, fullPath, handlerFunc)
}

// RegisterGroup registers a route group with the given name and builder function.
func RegisterGroup(groupName string, e *echo.Echo, builder func(g *echo.Group)) *echo.Echo {
	builder(e.Group(groupName))
	return e
}

// apiVersion is a middleware function that prefixes the request path with the version from the "version" header.
func apiVersion(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		version := req.Header.Get("version")
		if version != "" {
			req.URL.Path = fmt.Sprintf("/%s%s", version, req.URL.Path)
		}
		return next(c)
	}
}
