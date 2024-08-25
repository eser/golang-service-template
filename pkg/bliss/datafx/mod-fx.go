package datafx

import (
	"go.uber.org/fx"
)

var FxModule = fx.Module( //nolint:gochecknoglobals
	"data",
	fx.Provide(
		FxNew,
	),
)

type FxResult struct {
	fx.Out

	DataProvider DataProvider
}

func FxNew() (FxResult, error) {
	return FxResult{
		Out: fx.Out{},

		DataProvider: NewDataProvider(),
	}, nil
}
