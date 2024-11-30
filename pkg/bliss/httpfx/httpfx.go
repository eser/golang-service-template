package httpfx

import (
	"log/slog"

	"github.com/eser/go-service/pkg/bliss/di"
	"github.com/eser/go-service/pkg/bliss/metricsfx"
)

func RegisterDependencies(container di.Container, config *Config, mp metricsfx.MetricsProvider, logger *slog.Logger) {
	routes := NewRouter("/")
	httpService := NewHttpService(config, routes, mp, logger)

	di.RegisterFor[Router](container, routes)
	di.RegisterFor[HttpService](container, httpService)
	di.Register(container, httpService.InnerMetrics)
}
