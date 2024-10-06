package main

import (
	"context"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/http/echo/server"
	log "github.com/NekKkMirror/go-app/internal/pkg/logger"
)

func main() {
	logCfg := log.Config{
		LogLevel:      "info",
		LokiURL:       "",
		Labels:        map[string]string{"job": "test"},
		SendInterval:  time.Second,
		MaxBatchSize:  10,
		MaxBufferSize: 100,
		RetryCount:    3,
		RetryDelay:    500 * time.Millisecond,
	}

	logger := log.InitLogger(&logCfg)

	serverCfg := &server.EchoConfig{
		Host:               "localhost",
		Port:               "8080",
		BasePath:           "/api/v1",
		DebugErrorResponse: false,
		Timeout:            10 * time.Second,
		RateLimit:          100,
	}

	e := server.NewEchoServer(serverCfg, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := e.Run(ctx); err != nil {
		logger.Errorf("Error starting the server: %v", err)
	}
}
