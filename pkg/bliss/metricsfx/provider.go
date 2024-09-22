package metricsfx

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type MetricsProvider interface {
	GetRegistry() *prometheus.Registry
}

type MetricsProviderImpl struct {
	registry *prometheus.Registry
}

var _ MetricsProvider = (*MetricsProviderImpl)(nil)

func NewMetricsProvider() *MetricsProviderImpl {
	registry := prometheus.NewRegistry()

	// Register the Go collector (which collects runtime metrics)
	registry.MustRegister(collectors.NewGoCollector())

	// Register the process collector (which collects process-level metrics)
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return &MetricsProviderImpl{
		registry: registry,
	}
}

func (mp *MetricsProviderImpl) GetRegistry() *prometheus.Registry {
	return mp.registry
}
