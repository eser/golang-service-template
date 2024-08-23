package configfx

import (
	"errors"

	"go.uber.org/fx"
)

var ErrConfigDecoding = errors.New("config decoding error")

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
