package grpc

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

func (s *Server) shutdown() {
	s.Grpc.Stop()
	s.Grpc.GracefulStop()
}
