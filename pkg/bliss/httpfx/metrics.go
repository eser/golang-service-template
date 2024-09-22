package httpfx

import (
	"github.com/eser/go-service/pkg/bliss/metricsfx"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	mp metricsfx.MetricsProvider

	RequestsTotal *prometheus.CounterVec
}

func NewMetrics(mp metricsfx.MetricsProvider) *Metrics {
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	mp.GetRegistry().MustRegister(requestsTotal)

	return &Metrics{
		mp:            mp,
		RequestsTotal: requestsTotal,
	}
}
