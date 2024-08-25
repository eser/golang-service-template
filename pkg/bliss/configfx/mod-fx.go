package configfx

import (
	"go.uber.org/fx"
)

var FxModule = fx.Module( //nolint:gochecknoglobals
	"config",
	fx.Provide(
		FxNew,
	),
)

type FxResult struct {
	fx.Out

	ConfigLoader ConfigLoader
}

func FxNew() FxResult {
	return FxResult{
		Out: fx.Out{},

		ConfigLoader: NewConfigLoader(),
	}
}
