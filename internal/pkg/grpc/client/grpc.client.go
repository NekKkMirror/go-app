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

// NewGrpcClient creates a new gRPC client connection to the specified host and port.
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

func (g *grpcClient) Close() error {
	if g.conn == nil {
		return nil
	}
	return g.conn.Close()
}
