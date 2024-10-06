package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/logger"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	maxConnectionIdle = 5
	maxConnectionAge  = 5
	gRPCTime          = 10
	gRPCTimeout       = 15
)

type Config struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	Development bool   `mapstructure:"development"`
}

type Server struct {
	Grpc   *grpc.Server
	Config *Config
	Log    logger.ILogger
}

// NewGrpcServer creates a new gRPC server instance with the provided configuration and logger.
// It initializes the gRPC server with keepalive parameters and OpenTelemetry instrumentation.
//
// Parameters:
// - log: ILogger instance for logging.
// - config: Configuration for the gRPC server.
//
// Returns:
// - A pointer to a new Server instance.
func NewGrpcServer(log logger.ILogger, config *Config) *Server {
	serverOptions := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}

	s := grpc.NewServer(serverOptions...)

	return &Server{Grpc: s, Config: config, Log: log}
}

// RunGrpcServer starts the gRPC server with the provided context and optional configuration function.
// It listens on the specified host and port, applies keepalive parameters, and registers reflection service if in development mode.
// The server runs in a separate goroutine and shuts down gracefully when the provided context is done.
//
// Parameters:
// - ctx: Context for managing the lifecycle of the server.
// - configGrpc: Optional function to configure the gRPC server.
//
// Returns:
// - An error if the gRPC server fails to start or serve.
func (s *Server) RunGrpcServer(ctx context.Context, configGrpc ...func(grpcServer *grpc.Server)) error {
	address := net.JoinHostPort(s.Config.Host, s.Config.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "net.Listen failed to listen")
	}

	if len(configGrpc) > 0 {
		grpcFunc := configGrpc[0]
		if grpcFunc != nil {
			grpcFunc(s.Grpc)
		}
	}

	if s.Config.Development {
		reflection.Register(s.Grpc)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					s.Log.Infof("shutting down grpc server port: %s", s.Config.Port)
					s.shutdown()
					s.Log.Infof("grpc server exited properly")
					return
				}
			}
		}
	}()

	s.Log.Infof("grpc server is listen on port: %s", s.Config.Port)

	err = s.Grpc.Serve(listen)

	if err != nil {
		s.Log.Error(fmt.Sprintf("[grpcServer_RunGrpcServer.Serve] grpc server serve error: %+v", err))
	}

	return err
}

// shutdown gracefully stops the gRPC server and waits for all connections to be closed.
// It uses the Stop() method to stop accepting new connections and the GracefulStop() method to wait for all existing connections to be closed.
//
// This function is intended to be called when the server needs to be shut down gracefully.
// It should be called within a separate goroutine to avoid blocking the main execution flow.
//
// No parameters are required for this function.
//
// The function does not return any value.
func (s *Server) shutdown() {
	s.Grpc.Stop()
	s.Grpc.GracefulStop()
}
