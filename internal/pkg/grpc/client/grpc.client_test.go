package client

import (
	"testing"

	"github.com/NekKkMirror/go-app/internal/pkg/grpc/client/mocks"
	"github.com/NekKkMirror/go-app/internal/pkg/grpc/server"
	"github.com/stretchr/testify/assert"
	google_golang_orggrpc "google.golang.org/grpc"
)

func TestNewGrpcClient(t *testing.T) {
	config := &server.Config{
		Host: "localhost",
		Port: "50051",
	}

	client, err := NewGrpcClient(config)

	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestGrpcClient_GetGrpcConnection(t *testing.T) {
	mockClient := mocks.NewClient(t)
	expectedConn := &google_golang_orggrpc.ClientConn{}

	mockClient.On("GetGrpcConnection").Return(expectedConn)

	conn := mockClient.GetGrpcConnection()

	assert.Equal(t, expectedConn, conn)
	mockClient.AssertExpectations(t)
}

func TestGrpcClient_Close(t *testing.T) {
	mockClient := mocks.NewClient(t)

	mockClient.On("Close").Return(nil)

	err := mockClient.Close()

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}
