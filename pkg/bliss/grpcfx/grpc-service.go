package grpcfx

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/eser/go-service/pkg/bliss/metricsfx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcService interface {
	Server() *grpc.Server
	RegisterService(desc *grpc.ServiceDesc, impl any)
	Start(ctx context.Context) (func(), error)
}

type GrpcServiceImpl struct {
	InnerServer  *grpc.Server
	InnerMetrics *Metrics
	Config       *Config
	logger       *slog.Logger
}

var _ GrpcService = (*GrpcServiceImpl)(nil)

func NewGrpcService(config *Config, metricsProvider metricsfx.MetricsProvider, logger *slog.Logger) *GrpcServiceImpl {
	metrics := NewMetrics(metricsProvider)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggingInterceptor(logger),
			MetricsInterceptor(metrics),
		),
	)

	if config.Reflection {
		reflection.Register(server)
	}

	return &GrpcServiceImpl{
		InnerServer:  server,
		InnerMetrics: metrics,
		Config:       config,
		logger:       logger,
	}
}

func (gs *GrpcServiceImpl) Server() *grpc.Server {
	return gs.InnerServer
}

func (gs *GrpcServiceImpl) RegisterService(desc *grpc.ServiceDesc, impl any) {
	gs.InnerServer.RegisterService(desc, impl)
}

func (gs *GrpcServiceImpl) Start(ctx context.Context) (func(), error) {
	gs.logger.InfoContext(ctx, "GrpcService is starting...", slog.String("addr", gs.Config.Addr))

	listener, err := net.Listen("tcp", gs.Config.Addr)
	if err != nil {
		return nil, fmt.Errorf("GrpcService Net Listen error: %w", err)
	}

	go func() {
		if err := gs.InnerServer.Serve(listener); err != nil {
			gs.logger.ErrorContext(ctx, "GrpcService Serve error", slog.Any("error", err))
		}
	}()

	cleanup := func() {
		gs.logger.InfoContext(ctx, "Shutting down gRPC server...")

		stopped := make(chan struct{})
		go func() {
			gs.InnerServer.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			gs.logger.InfoContext(ctx, "GrpcService has gracefully stopped.")
		case <-time.After(gs.Config.GracefulShutdownTimeout):
			gs.logger.WarnContext(ctx, "GrpcService shutdown timeout exceeded, forcing stop")
			gs.InnerServer.Stop()
		}
	}

	return cleanup, nil
}
