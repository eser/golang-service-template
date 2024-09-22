package metricsfx

import (
	"go.uber.org/fx"
)

var FxModule = fx.Module( //nolint:gochecknoglobals
	"metrics",
	fx.Provide(
		FxNew,
	),
)

type FxResult struct {
	fx.Out

	MetricsProvider MetricsProvider
}

func FxNew() FxResult {
	return FxResult{
		Out: fx.Out{},

		MetricsProvider: NewMetricsProvider(),
	}
}
