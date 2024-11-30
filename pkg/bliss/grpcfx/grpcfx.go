package grpcfx

import (
	"log/slog"

	"github.com/eser/go-service/pkg/bliss/di"
	"github.com/eser/go-service/pkg/bliss/metricsfx"
)

func RegisterDependencies(container di.Container, config *Config, mp metricsfx.MetricsProvider, logger *slog.Logger) {
	service := NewGrpcService(config, mp, logger)

	di.RegisterFor[GrpcService](container, service)
	di.Register(container, service.InnerMetrics)
}
