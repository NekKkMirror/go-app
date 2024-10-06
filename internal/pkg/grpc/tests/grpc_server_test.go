package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/grpc"
	"github.com/NekKkMirror/go-app/internal/pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRunGrpcServer(t *testing.T) {
	mockLogger := mocks.NewILogger(t)
	config := &grpc.Config{
		Host:        "localhost",
		Port:        "50051",
		Development: true,
	}

	server := grpc.NewGrpcServer(mockLogger, config)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	mockLogger.On("Infof", "grpc server is listen on port: %s", config.Port).Return(nil)
	mockLogger.On("Infof", "shutting down grpc server port: %s", config.Port).Return(nil)
	mockLogger.On("Infof", "grpc server exited properly").Return(nil)

	err := server.RunGrpcServer(ctx)

	assert.NoError(t, err)

	mockLogger.AssertExpectations(t)
}
