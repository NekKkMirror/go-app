package server

import (
	"context"
	"testing"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRunGrpcServer(t *testing.T) {
	mockLogger := mocks.NewILogger(t)
	config := &Config{
		Host:        "localhost",
		Port:        "50051",
		Development: true,
	}

	server := NewGrpcServer(mockLogger, config)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	mockLogger.On("Infof", "gRPC server listening on port: %s", config.Port).Return(nil)
	mockLogger.On("Infof", "shutting down gRPC server on port: %s", config.Port).Return(nil)
	mockLogger.On("Infof", "gRPC server exited properly").Return(nil)

	err := server.RunGrpcServer(ctx)

	assert.NoError(t, err)

	mockLogger.AssertExpectations(t)
}
