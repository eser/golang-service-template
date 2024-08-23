package datafx

import (
	"go.uber.org/fx"
)

var FxModule = fx.Module( //nolint:gochecknoglobals
	"data",
	fx.Provide(
		New,
	),
)

type FxResult struct {
	fx.Out

	DataProvider DataProvider
}

func New() (FxResult, error) {
	return FxResult{ //nolint:exhaustruct
		DataProvider: NewDataProvider(),
	}, nil
}
