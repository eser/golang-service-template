package configfx

import (
	"errors"

	"go.uber.org/fx"
)

var ErrConfigDecoding = errors.New("config decoding error")

var Module = fx.Module( //nolint:gochecknoglobals
	"config",
	fx.Provide(
		New,
	),
)

type Result struct {
	fx.Out

	ConfigLoader ConfigLoader
}

func New() (Result, error) {
	return Result{ //nolint:exhaustruct
		ConfigLoader: NewConfigLoader(),
	}, nil
}
