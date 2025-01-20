package metricsfx

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type MetricsProvider struct {
	registry *prometheus.Registry
}

func NewMetricsProvider() *MetricsProvider {
	registry := prometheus.NewRegistry()

	// Register the Go collector (which collects runtime metrics)
	registry.MustRegister(collectors.NewGoCollector())

	// Register the process collector (which collects process-level metrics)
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})) //nolint:exhaustruct

	return &MetricsProvider{
		registry: registry,
	}
}

func (mp *MetricsProvider) GetRegistry() *prometheus.Registry {
	return mp.registry
}
