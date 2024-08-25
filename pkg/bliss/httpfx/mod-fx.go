package httpfx

import (
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
}

func FxNew(config *Config) FxResult {
	routes := NewRouter("/")
	httpService := NewHttpService(config, routes)

	return FxResult{
		Out: fx.Out{},

		HttpService: httpService,
		Routes:      routes,
	}
}

func registerHooks(lc fx.Lifecycle, hs HttpService) {
	lc.Append(fx.Hook{
		OnStart: hs.Start,
		OnStop:  hs.Stop,
	})
}
