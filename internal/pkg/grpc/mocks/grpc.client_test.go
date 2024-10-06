package mocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestGrpcClientMock_Close(t *testing.T) {
	client := newGrpcClient(t)

	client.On("Close").Return(nil)

	err := client.Close()

	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestGrpcClientMock_GetGrpcConnection(t *testing.T) {
	client := newGrpcClient(t)
	expectedConn := &grpc.ClientConn{}

	client.On("GetGrpcConnection").Return(expectedConn)

	conn := client.GetGrpcConnection()

	assert.Equal(t, expectedConn, conn)
	client.AssertExpectations(t)
}
