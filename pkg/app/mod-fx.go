package app

import (
	"fmt"

	"github.com/eser/go-service/pkg/app/config"
	"github.com/eser/go-service/pkg/app/routes/home"
	"github.com/eser/go-service/pkg/app/routes/protected"
	"github.com/eser/go-service/pkg/bliss"
	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/httpfx/middlewares"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/healthcheck"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/openapi"
	"go.uber.org/fx"
)

var FxModule = fx.Module( //nolint:gochecknoglobals
	"app",
	fx.Invoke(
		RegisterMiddlewares,
		home.IndexRoutes,
		protected.IndexRoutes,
	),
	fx.Provide(
		bliss.LoadConfig[config.AppConfig](LoadConfig),
	),
	healthcheck.FxModule,
	openapi.FxModule,
)

func LoadConfig(cl configfx.ConfigLoader) (*config.AppConfig, error) {
	appConfig := &config.AppConfig{} //nolint:exhaustruct

	err := cl.Load(
		appConfig,

		cl.FromJsonFile("config.json"),
		cl.FromEnvFile(".env"),
		cl.FromSystemEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return appConfig, nil
}

func RegisterMiddlewares(routes httpfx.Router, httpMetrics *httpfx.Metrics, appConfig *config.AppConfig) {
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResolveAddressMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.CorrelationIdMiddleware())
	routes.Use(middlewares.CorsMiddleware())
	routes.Use(middlewares.MetricsMiddleware(httpMetrics))
}
