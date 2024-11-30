package grpcfx

import (
	"github.com/eser/go-service/pkg/bliss/metricsfx"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	mp metricsfx.MetricsProvider

	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

func NewMetrics(metricsProvider metricsfx.MetricsProvider) *Metrics {
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{ //nolint:exhaustruct
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "code"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{ //nolint:exhaustruct
			Name: "grpc_request_duration_seconds",
			Help: "gRPC request duration in seconds",
		},
		[]string{"method"},
	)

	metricsProvider.GetRegistry().MustRegister(requestsTotal, requestDuration)

	return &Metrics{
		mp:              metricsProvider,
		RequestsTotal:   requestsTotal,
		RequestDuration: requestDuration,
	}
}
