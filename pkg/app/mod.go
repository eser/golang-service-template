package app

import (
	"net/http"

	"github.com/eser/go-service/pkg/bliss"
	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/httpfx/middlewares"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/healthcheck"
	"github.com/eser/go-service/pkg/bliss/httpfx/modules/openapi"
	"go.uber.org/fx"
)

var appConfig = AppConfig{}

var Module = fx.Module( //nolint:gochecknoglobals
	"app",
	fx.Invoke(
		LoadConfig,
		RegisterRoutes,
	),
	healthcheck.Module,
	openapi.Module,
)

func LoadConfig(conf configfx.ConfigLoader) {
	conf.Load(
		&appConfig,

		conf.FromJsonFile("config.json"),
		conf.FromEnvFile(".env"),
		conf.FromSystemEnv(),
	)
}

func RegisterRoutes(routes httpfx.Router) {
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.CorrelationIdMiddleware())
	routes.Use(middlewares.CorsMiddleware())

	routes.
		Route("GET /", func(ctx *httpfx.Context) httpfx.Response {
			return ctx.Results.PlainText("Hello, World!")
		}).
		HasSummary("Homepage").
		HasDescription("This is the homepage of the service.").
		HasResponse(http.StatusOK)
}

func New() *fx.App {
	return fx.New(
		// fx.WithLogger(bliss.GetFxLogger),
		bliss.Module,
		Module,
	)
}
