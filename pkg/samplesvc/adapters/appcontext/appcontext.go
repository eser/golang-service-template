package appcontext

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/eser/ajan/configfx"
	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/logfx"
	"github.com/eser/ajan/metricsfx"
)

var ErrInitFailed = errors.New("failed to initialize app context")

type AppContext struct {
	Config  *AppConfig
	Logger  *slog.Logger
	Metrics *metricsfx.MetricsProvider
	Data    *datafx.Registry
}

func NewAppContext(ctx context.Context) (*AppContext, error) {
	appContext := &AppContext{} //nolint:exhaustruct

	// config
	cl := configfx.NewConfigManager()

	appContext.Config = &AppConfig{} //nolint:exhaustruct

	err := cl.LoadDefaults(appContext.Config)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// logger
	appContext.Logger, err = logfx.NewLoggerAsDefault(os.Stdout, &appContext.Config.Log)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// metrics
	appContext.Metrics = metricsfx.NewMetricsProvider()

	err = appContext.Metrics.RegisterNativeCollectors()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	// data
	appContext.Data = datafx.NewRegistry(appContext.Logger)

	err = appContext.Data.LoadFromConfig(ctx, &appContext.Config.Data)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInitFailed, err)
	}

	return appContext, nil
}
