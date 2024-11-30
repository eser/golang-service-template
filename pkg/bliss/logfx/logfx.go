package logfx

import (
	"github.com/eser/go-service/pkg/bliss/di"
)

func RegisterDependencies(container di.Container, config *Config) error {
	logger, err := NewLoggerAsDefault(config)
	if err != nil {
		return err
	}

	di.Register(container, logger)

	return nil
}
