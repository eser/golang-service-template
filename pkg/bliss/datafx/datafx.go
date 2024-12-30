package datafx

import (
	"fmt"
	"log/slog"

	"github.com/eser/go-service/pkg/bliss/di"
	_ "github.com/lib/pq"
)

func RegisterDependencies(container di.Container, config *Config, logger *slog.Logger) error {
	dp := NewDataManager(logger)

	di.RegisterFor[DataProvider](container, dp)

	err := dp.LoadFromConfig(config)
	if err != nil {
		return fmt.Errorf("failed to load database connections from config: %w", err)
	}

	return nil
}
