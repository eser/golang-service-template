package httpfx

import (
	"github.com/eser/go-service/pkg/bliss/metricsfx"
	"go.uber.org/fx"
)

var FxModule = fx.Module( //nolint:gochecknoglobals
	"httpservice",
	fx.Provide(
		FxNew,
	),
	fx.Invoke(
		registerHooks,
	),
)

type FxResult struct {
	fx.Out

	HttpService HttpService
	Routes      Router
	Metrics     *Metrics
}

func FxNew(config *Config, mp metricsfx.MetricsProvider) FxResult {
	routes := NewRouter("/")
	httpService := NewHttpService(config, routes)
	metrics := NewMetrics(mp)

	return FxResult{
		Out: fx.Out{},

		HttpService: httpService,
		Routes:      routes,
		Metrics:     metrics,
	}
}

func registerHooks(lc fx.Lifecycle, hs HttpService) {
	lc.Append(fx.Hook{
		OnStart: hs.Start,
		OnStop:  hs.Stop,
	})
}
