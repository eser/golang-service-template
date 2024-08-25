package logfx

import (
	"log/slog"
	"os"
)

func NewLogger(config *Config) (*slog.Logger, error) {
	handler, err := NewHandler(os.Stderr, config)
	if err != nil {
		return nil, err
	}

	return slog.New(handler), nil
}

func NewLoggerAsDefault(config *Config) (*slog.Logger, error) {
	logger, err := NewLogger(config)
	if err != nil {
		return nil, err
	}

	slog.SetDefault(logger)

	return logger, nil
}
