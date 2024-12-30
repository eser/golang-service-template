package samplehttp

import (
	"context"
	"fmt"

	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/bliss/di"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/httpfx/middlewares"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/healthcheck"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/openapi"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/profiling"
	"github.com/eser/go-service/pkg/bliss/lib"
	"github.com/eser/go-service/pkg/bliss/logfx"
	"github.com/eser/go-service/pkg/bliss/metricsfx"
)

func LoadConfig(loader configfx.ConfigLoader) (*AppConfig, *logfx.Config, *httpfx.Config, *datafx.Config, error) {
	appConfig := &AppConfig{} //nolint:exhaustruct

	err := loader.LoadDefaults(appConfig)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	return appConfig, &appConfig.Log, &appConfig.Http, &appConfig.Data, nil
}

func RegisterHttpMiddlewares(routes httpfx.Router, httpMetrics *httpfx.Metrics, appConfig *AppConfig) error {
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

	appConfig := &AppConfig{} //nolint:exhaustruct

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
	dataProvider := datafx.NewDataManager(logger)
	err = dataProvider.LoadFromConfig(&appConfig.Data)
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
	RegisterHttpRoutes(routes, appConfig, logger, dataProvider)

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

func RunWithDi() error {
	err := di.RegisterFn(
		di.Default,
		configfx.RegisterDependencies,
		LoadConfig,

		logfx.RegisterDependencies,
		metricsfx.RegisterDependencies,
		httpfx.RegisterDependencies,
		datafx.RegisterDependencies,

		RegisterHttpMiddlewares,

		healthcheck.RegisterHttpRoutes,
		openapi.RegisterHttpRoutes,
		profiling.RegisterHttpRoutes,

		RegisterHttpRoutes,
	)
	if err != nil {
		panic(err)
	}

	run := di.CreateInvoker(
		di.Default,
		func(
			httpService httpfx.HttpService,
		) error {
			ctx := context.Background()

			cleanup, err := httpService.Start(ctx)
			if err != nil {
				return err //nolint:wrapcheck
			}

			lib.WaitForSignal()

			cleanup()

			return nil
		},
	)

	di.Seal(di.Default)

	err = run()

	return err
}
