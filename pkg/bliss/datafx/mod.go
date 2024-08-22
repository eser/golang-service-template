package datafx

import (
	"go.uber.org/fx"
)

var Module = fx.Module( //nolint:gochecknoglobals
	"data",
	fx.Provide(
		New,
	),
)

type Result struct {
	fx.Out

	DataProvider DataProvider
}

func New() (Result, error) {
	return Result{ //nolint:exhaustruct
		DataProvider: NewDataProvider(),
	}, nil
}
