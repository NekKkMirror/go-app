package client

import "google.golang.org/grpc"

// Client interface defines methods for working with gRPC client connections.
//
//go:generate mockery --name Client
type Client interface {
	// GetGrpcConnection returns the gRPC client connection.
	GetGrpcConnection() *grpc.ClientConn

	// Close closes the gRPC client connection and releases all associated resources.
	Close() error
}
