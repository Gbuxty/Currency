package server

import (
	"GetCurrency/internal/config"
	"GetCurrency/internal/transport/grpc/handlers"
	"GetCurrency/pkg/logger"
	rate_pb "GetCurrency/proto/rate.pb"
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	server   *grpc.Server
	logger   logger.Logger
	cfg      *config.ServerConfig
	handlers *handlers.RateHandler
}

func NewGrpcServer(logger logger.Logger, cfg *config.ServerConfig, handlers *handlers.RateHandler) *GrpcServer {
	return &GrpcServer{
		server:   grpc.NewServer(),
		logger:   logger,
		handlers: handlers,
		cfg:      cfg,
	}
}

func (s *GrpcServer) Start(ctx context.Context) error {
	
	addr := fmt.Sprintf(":%d", s.cfg.Port)

	lc := &net.ListenConfig{}

	lis, err := lc.Listen(ctx, "tcp", addr)

	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	rate_pb.RegisterRateServiceServer(s.server, s.handlers)

	go func() {
		if err := s.server.Serve(lis); err != nil {
			s.logger.Errorf("failed to serve gRPC server: %v", err)
		}
	}()

	s.logger.Infof("gRPC server started port: %s", addr)
	return nil
}

func (s *GrpcServer) Stop() {
	s.logger.Infof("Stopping gRPC server...")
	s.server.GracefulStop()
}
