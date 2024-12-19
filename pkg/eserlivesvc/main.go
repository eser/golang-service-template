package eserlivesvc

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
	"github.com/eser/go-service/pkg/eserlivesvc/entities/users"
	"github.com/eser/go-service/pkg/eserlivesvc/shared"
)

func LoadConfig(loader configfx.ConfigLoader) (*shared.AppConfig, *logfx.Config, *httpfx.Config, *datafx.Config, error) { //nolint:lll
	appConfig := &shared.AppConfig{} //nolint:exhaustruct

	err := loader.LoadDefaults(appConfig)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	return appConfig, &appConfig.Log, &appConfig.Http, &appConfig.Data, nil
}

func RegisterHttpMiddlewares(routes httpfx.Router, httpMetrics *httpfx.Metrics, appConfig *shared.AppConfig) error {
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResolveAddressMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.CorrelationIdMiddleware())
	routes.Use(middlewares.CorsMiddleware())
	routes.Use(middlewares.MetricsMiddleware(httpMetrics))

	return nil
}

func Run() error {
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

		users.RegisterHttpRoutes,
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
