package server

import (
	"context"
	"fmt"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/logger"
	"github.com/labstack/echo/v4"
)

const (
	MaxHeaderBytes = 1 << 20
	ReadTimeout    = 15 * time.Second
	WriteTimeout   = 15 * time.Second
)

type EchoConfig struct {
	Host               string   `mapstructure:"host"`
	Port               string   `mapstructure:"port" validate:"required"`
	Development        string   `mapstructure:"development"`
	BasePath           string   `mapstructure:"basePath" validate:"required"`
	DebugErrorResponse bool     `mapstructure:"debugErrorResponse"`
	IgnoreLogUrls      []string `mapstructure:"ignoreLogUrls"`
	Timeout            string   `mapstructure:"timeout"`
}

func NewEchoServer() *echo.Echo {
	e := echo.New()
	return e
}

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

	err := echo.Start(cfg.Port)

	return err
}

func ApplyVersioningFromHeader(echo *echo.Echo) {
	echo.Pre(apiVersion)
}

func RegisterGroupFunc(groupName string, echo *echo.Echo, builder func(g *echo.Group)) *echo.Echo {
	builder(echo.Group(groupName))

	return echo
}

func apiVersion(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		headers := req.Header

		apiVersion := headers.Get("version")

		req.URL.Path = fmt.Sprintf("%s%s", apiVersion, req.URL.Path)

		return next(c)
	}
}
