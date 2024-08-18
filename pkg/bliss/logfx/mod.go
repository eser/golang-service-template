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
	var level slog.Level

	err := level.UnmarshalText([]byte(config.Level))
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: replaceAttr,
		AddSource:   config.AddSource,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger, nil
}
