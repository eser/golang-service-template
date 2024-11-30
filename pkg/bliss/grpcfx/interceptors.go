package grpcfx

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)

		// Log in same format as httpfx
		logger.InfoContext(ctx, "RPC Call",
			slog.String("method", info.FullMethod),
			slog.String("duration", duration.String()),
		)

		return resp, err
	}
}

func MetricsInterceptor(metrics *Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		st, _ := status.FromError(err)

		metrics.RequestsTotal.WithLabelValues(
			info.FullMethod,
			st.Code().String(),
		).Inc()

		metrics.RequestDuration.WithLabelValues(
			info.FullMethod,
		).Observe(duration.Seconds())

		return resp, err
	}
}
