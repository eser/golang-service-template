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

	return &MetricsProvider{
		registry: registry,
	}
}

func (mp *MetricsProvider) RegisterNativeCollectors() error {
	// Register the Go collector (which collects runtime metrics)
	goCollector := collectors.NewGoCollector()
	if err := mp.registry.Register(goCollector); err != nil {
		return err //nolint:wrapcheck
	}

	// Register the process collector (which collects process-level metrics)
	processCollector := collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
		PidFn:        nil,
		Namespace:    "", // TODO(@eser) fill with namespace information from deployment
		ReportErrors: false,
	})
	if err := mp.registry.Register(processCollector); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}

func (mp *MetricsProvider) GetRegistry() *prometheus.Registry {
	return mp.registry
}
