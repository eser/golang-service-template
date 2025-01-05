package samplesvc

import (
	"context"
	"fmt"

	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/httpfx/middlewares"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/healthcheck"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/openapi"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/profiling"
	"github.com/eser/go-service/pkg/bliss/lib"
	"github.com/eser/go-service/pkg/bliss/logfx"
	"github.com/eser/go-service/pkg/bliss/metricsfx"
	"github.com/eser/go-service/pkg/samplesvc/adapters/config"
	"github.com/eser/go-service/pkg/samplesvc/adapters/http"
)

func LoadConfig(loader configfx.ConfigLoader) (*config.AppConfig, *logfx.Config, *httpfx.Config, *datafx.Config, error) { //nolint:lll
	appConfig := &config.AppConfig{} //nolint:exhaustruct

	err := loader.LoadDefaults(appConfig)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	return appConfig, &appConfig.Log, &appConfig.Http, &appConfig.Data, nil
}

func RegisterHttpMiddlewares(routes httpfx.Router, httpMetrics *httpfx.Metrics, appConfig *config.AppConfig) error {
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResolveAddressMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.CorrelationIdMiddleware())
	routes.Use(middlewares.CorsMiddleware())
	routes.Use(middlewares.MetricsMiddleware(httpMetrics))

	return nil
}

func Run() error {
	// config
	cl := configfx.NewConfigManager()

	appConfig := &config.AppConfig{} //nolint:exhaustruct

	err := cl.LoadDefaults(appConfig)
	if err != nil {
		return err //nolint:wrapcheck
	}

	// logger
	logger, err := logfx.NewLoggerAsDefault(&appConfig.Log)
	if err != nil {
		return err //nolint:wrapcheck
	}

	// metrics
	mp := metricsfx.NewMetricsProvider()

	// http service
	routes := httpfx.NewRouter("/")
	httpService := httpfx.NewHttpService(&appConfig.Http, routes, mp, logger)

	// data
	dataRegistry := datafx.NewRegistry(logger)

	err = dataRegistry.LoadFromConfig(context.TODO(), &appConfig.Data)
	if err != nil {
		return err //nolint:wrapcheck
	}

	// http middlewares
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResolveAddressMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.CorrelationIdMiddleware())
	routes.Use(middlewares.CorsMiddleware())
	routes.Use(middlewares.MetricsMiddleware(httpService.InnerMetrics))

	// http modules
	healthcheck.RegisterHttpRoutes(routes, &appConfig.Http)
	openapi.RegisterHttpRoutes(routes, &appConfig.Http)
	profiling.RegisterHttpRoutes(routes, &appConfig.Http)

	// http routes
	http.RegisterHttpRoutes(routes, appConfig, logger, dataRegistry)

	// run
	ctx := context.Background()

	cleanup, err := httpService.Start(ctx)
	if err != nil {
		return err //nolint:wrapcheck
	}

	defer cleanup()

	lib.WaitForSignal()

	return nil
}
