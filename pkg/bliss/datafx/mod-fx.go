package datafx

import (
	"database/sql"

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

func FxNew(db *sql.DB) FxResult {
	return FxResult{
		Out: fx.Out{},

		DataProvider: NewDataProvider(db),
	}
}
