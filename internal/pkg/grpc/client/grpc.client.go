package client

import (
	"fmt"

	grpc2 "github.com/NekKkMirror/go-app/internal/pkg/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClient struct {
	conn *grpc.ClientConn
}

//go:generate mockery --name Client
type Client interface {
	GetGrpcConnection() *grpc.ClientConn
	Close() error
}

// NewGrpcClient creates a new gRPC client connection to the specified host and port.
// It uses the provided configuration to establish the connection.
//
// The function takes a pointer to a grpc2.Config struct as a parameter.
// The Config struct contains the host and port information for the gRPC server.
//
// The function returns a Client interface and an error.
// If the connection is successfully established, the Client interface will be implemented by the grpcClient struct,
// and the error will be nil.
// If an error occurs during the connection establishment, the Client interface will be nil,
// and the error will contain the details of the failure.
func NewGrpcClient(config *grpc2.Config) (Client, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", config.Host, config.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClient{conn: conn}, nil
}

// GetGrpcConnection returns the gRPC client connection.
// It retrieves and returns the underlying grpc.ClientConn instance.
func (g *grpcClient) GetGrpcConnection() *grpc.ClientConn {
	return g.conn
}

// Close closes the gRPC client connection.
// It releases all resources associated with the connection.
//
// The function returns an error if the connection cannot be closed successfully.
// If the connection is already closed, the function will return nil.
func (g *grpcClient) Close() error {
	return g.conn.Close()
}
