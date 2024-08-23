package configfx

import (
	"go.uber.org/fx"
)

var FxModule = fx.Module( //nolint:gochecknoglobals
	"config",
	fx.Provide(
		New,
	),
)

type FxResult struct {
	fx.Out

	ConfigLoader ConfigLoader
}

func New() (FxResult, error) {
	return FxResult{ //nolint:exhaustruct
		ConfigLoader: NewConfigLoader(),
	}, nil
}
