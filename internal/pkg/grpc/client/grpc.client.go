package client

import (
	"fmt"

	grpc2 "github.com/NekKkMirror/go-app/internal/pkg/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client interface defines methods for working with gRPC client connections.
//
//go:generate mockery --name Client
type Client interface {
	GetGrpcConnection() *grpc.ClientConn
	Close() error
}

type grpcClient struct {
	conn *grpc.ClientConn
}

// NewGrpcClient creates a new gRPC client connection to the specified host and port.
//
// The function takes a pointer to a grpc2.Config struct as a parameter, which contains
// the host and port information for the gRPC server.
//
// It returns a Client interface and an error.
// If the connection is successfully established, the Client interface will be
// implemented by the grpcClient struct, and the error will be nil.
// If an error occurs during the connection establishment, the Client interface will be nil,
// and the error will contain the details of the failure.
func NewGrpcClient(config *grpc2.Config) (Client, error) {
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %w", address, err)
	}

	return &grpcClient{conn: conn}, nil
}

// GetGrpcConnection returns the gRPC client connection.
func (g *grpcClient) GetGrpcConnection() *grpc.ClientConn {
	return g.conn
}

// Close closes the gRPC client connection and releases all associated resources.
//
// It returns an error if the connection cannot be closed successfully.
// If the connection is already closed, it returns nil.
func (g *grpcClient) Close() error {
	if g.conn == nil {
		return nil
	}
	return g.conn.Close()
}
