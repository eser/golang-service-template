package datafx

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/eser/go-service/pkg/bliss/di"
	_ "github.com/lib/pq"
)

func RegisterDependencies(container di.Container, config *Config, logger *slog.Logger) error {
	dr := NewRegistry(logger)

	di.RegisterFor[*Registry](container, dr)

	err := dr.LoadFromConfig(context.TODO(), config)
	if err != nil {
		return fmt.Errorf("failed to load database connections from config: %w", err)
	}

	return nil
}
