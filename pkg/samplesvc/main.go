package samplesvc

import (
	"context"
	"fmt"
	"os"

	"github.com/eser/ajan/configfx"
	"github.com/eser/ajan/datafx"
	"github.com/eser/ajan/httpfx"
	"github.com/eser/ajan/httpfx/middlewares"
	"github.com/eser/ajan/httpfx/modules/healthcheck"
	"github.com/eser/ajan/httpfx/modules/openapi"
	"github.com/eser/ajan/httpfx/modules/profiling"
	"github.com/eser/ajan/lib"
	"github.com/eser/ajan/logfx"
	"github.com/eser/ajan/metricsfx"
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
	ctx := context.Background()

	// config
	cl := configfx.NewConfigManager()

	appConfig := &config.AppConfig{} //nolint:exhaustruct

	err := cl.LoadDefaults(appConfig)
	if err != nil {
		return err //nolint:wrapcheck
	}

	// logger
	logger, err := logfx.NewLoggerAsDefault(os.Stdout, &appConfig.Log)
	if err != nil {
		return err //nolint:wrapcheck
	}

	// metrics
	mp := metricsfx.NewMetricsProvider() //nolint:varnamelen

	err = mp.RegisterNativeCollectors()
	if err != nil {
		return err //nolint:wrapcheck
	}

	// http service
	routes := httpfx.NewRouter("/")
	httpService := httpfx.NewHttpService(&appConfig.Http, routes, mp, logger)

	// data
	dataRegistry := datafx.NewRegistry(logger)

	err = dataRegistry.LoadFromConfig(ctx, &appConfig.Data)
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
	cleanup, err := httpService.Start(ctx)
	if err != nil {
		return err //nolint:wrapcheck
	}

	defer cleanup()

	lib.WaitForSignal()

	return nil
}
