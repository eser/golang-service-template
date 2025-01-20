package logfx

import (
	"io"
	"log/slog"
)

func NewLogger(w io.Writer, config *Config) (*slog.Logger, error) {
	handler, err := NewHandler(w, config)
	if err != nil {
		return nil, err
	}

	return slog.New(handler), nil
}

func NewLoggerAsDefault(w io.Writer, config *Config) (*slog.Logger, error) {
	logger, err := NewLogger(w, config)
	if err != nil {
		return nil, err
	}

	slog.SetDefault(logger)

	return logger, nil
}
