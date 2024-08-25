package logfx

import (
	"log/slog"

	"go.uber.org/fx"
)

var FxModule = fx.Module( //nolint:gochecknoglobals
	"log",
	fx.Provide(
		FxNew,
	),
)

type FxResult struct {
	fx.Out

	Logger *slog.Logger
}

func FxNew(config *Config) (FxResult, error) {
	logger, err := NewLogger(config)
	if err != nil {
		return FxResult{}, err
	}

	return FxResult{
		Out: fx.Out{},

		Logger: logger,
	}, nil
}
