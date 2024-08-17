package app

import (
	"fmt"
	"net/http"

	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/httpfx/middlewares"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/healthcheck"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/openapi"
	"go.uber.org/fx"
)

var Module = fx.Module( //nolint:gochecknoglobals
	"app",
	fx.Invoke(
		RegisterRoutes,
	),
	fx.Provide(
		LoadConfig,
		RegisterHttpConfig,
	),
	healthcheck.Module,
	openapi.Module,
)

func LoadConfig(conf configfx.ConfigLoader) (*AppConfig, error) {
	appConfig := &AppConfig{} //nolint:exhaustruct

	err := conf.Load(
		appConfig,

		conf.FromJsonFile("config.json"),
		conf.FromEnvFile(".env"),
		conf.FromSystemEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return appConfig, nil
}

func RegisterHttpConfig(appConfig *AppConfig) *httpfx.Config {
	return &appConfig.Http
}

func RegisterRoutes(routes httpfx.Router, appConfig *AppConfig) {
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.CorrelationIdMiddleware())
	routes.Use(middlewares.CorsMiddleware())

	routes.
		Route("GET /", func(ctx *httpfx.Context) httpfx.Response {
			message := fmt.Sprintf("Hello from %s!", appConfig.AppName)

			return ctx.Results.PlainText(message)
		}).
		HasSummary("Homepage").
		HasDescription("This is the homepage of the service.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /protected", middlewares.AuthMiddleware(), func(ctx *httpfx.Context) httpfx.Response {
			message := fmt.Sprintf("Hello from %s! this endpoint is protected!", appConfig.AppName)

			return ctx.Results.PlainText(message)
		}).
		HasSummary("Protected page").
		HasDescription("A page protected with JWT auth.").
		HasResponse(http.StatusOK)
}
