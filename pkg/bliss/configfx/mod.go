package configfx

import (
	"errors"

	"go.uber.org/fx"
)

var ErrConfigDecoding = errors.New("config decoding error")

//nolint:gochecknoglobals
var Module = fx.Module(
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
	return Result{
		ConfigLoader: NewConfigLoader(),
	}, nil
}
