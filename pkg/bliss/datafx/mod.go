package datafx

import (
	"go.uber.org/fx"
)

//nolint:gochecknoglobals
var Module = fx.Module(
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
