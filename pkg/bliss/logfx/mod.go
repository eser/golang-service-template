package logfx

import (
	"log/slog"
	"os"

	"go.uber.org/fx"
)

var FxModule = fx.Module( //nolint:gochecknoglobals
	"log",
	fx.Provide(
		RegisterLogger,
	),
)

func RegisterLogger(config *Config) (*slog.Logger, error) {
	handler, err := NewHandler(os.Stderr, config)
	if err != nil {
		return nil, err
	}

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger, nil
}
