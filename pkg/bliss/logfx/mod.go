package logfx

import (
	"log/slog"
	"os"

	"go.uber.org/fx"
)

//nolint:gochecknoglobals
var Module = fx.Module(
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
