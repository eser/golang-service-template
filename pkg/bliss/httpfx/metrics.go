package httpfx

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	mp MetricsProvider

	RequestsTotal *prometheus.CounterVec
}

func NewMetrics(mp MetricsProvider) *Metrics { //nolint:varnamelen
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{ //nolint:exhaustruct
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
