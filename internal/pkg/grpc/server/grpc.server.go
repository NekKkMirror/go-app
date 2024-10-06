package server

import (
	"context"
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
	maxConnectionIdle = 5 * time.Minute
	maxConnectionAge  = 5 * time.Minute
	gRPCTime          = 10 * time.Minute
	gRPCTimeout       = 15 * time.Second
)

// Config contains the configuration for the gRPC server.
type Config struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	Development bool   `mapstructure:"development"`
}

// Server wraps the gRPC server along with its configuration and logger.
type Server struct {
	Grpc   *grpc.Server
	Config *Config
	Log    logger.ILogger
}

// NewGrpcServer creates a new gRPC server instance with the provided configuration and logger.
//
// It initializes the gRPC server with keepalive parameters and OpenTelemetry instrumentation.
func NewGrpcServer(log logger.ILogger, config *Config) *Server {
	serverOptions := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle,
			Timeout:           gRPCTimeout,
			MaxConnectionAge:  maxConnectionAge,
			Time:              gRPCTime,
		}),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}

	s := grpc.NewServer(serverOptions...)

	return &Server{Grpc: s, Config: config, Log: log}
}

// RunGrpcServer starts the gRPC server and listens on the specified host and port.
//
// The server runs in a separate goroutine and shuts down gracefully when the provided context is done.
func (s *Server) RunGrpcServer(ctx context.Context, configGrpc ...func(grpcServer *grpc.Server)) error {
	address := net.JoinHostPort(s.Config.Host, s.Config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "net.Listen failed to listen")
	}

	// Apply additional server configurations if provided
	if len(configGrpc) > 0 && configGrpc[0] != nil {
		configGrpc[0](s.Grpc)
	}

	// Register reflection for development mode
	if s.Config.Development {
		reflection.Register(s.Grpc)
	}

	go s.handleServerShutdown(ctx)

	s.Log.Infof("gRPC server listening on port: %s", s.Config.Port)

	if err := s.Grpc.Serve(listener); err != nil {
		s.Log.Errorf("gRPC server serve error: %v", err)
		return err
	}

	return nil
}

// handleServerShutdown listens for context cancellation to shutdown the server gracefully.
func (s *Server) handleServerShutdown(ctx context.Context) {
	<-ctx.Done()
	s.Log.Infof("shutting down gRPC server on port: %s", s.Config.Port)

	s.Grpc.GracefulStop()
	s.Log.Infof("gRPC server exited properly")
}
