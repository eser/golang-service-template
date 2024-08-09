package app

import (
	"net/http"

	"github.com/eser/go-service/pkg/bliss"
	"github.com/eser/go-service/pkg/bliss/httpservice"
	"github.com/eser/go-service/pkg/bliss/httpservice/middlewares"
	"github.com/eser/go-service/pkg/bliss/httpservice/modules/healthcheck"
	"github.com/eser/go-service/pkg/bliss/httpservice/modules/openapi"
	"go.uber.org/fx"
)

var Module = fx.Module( //nolint:gochecknoglobals
	"app",
	fx.Invoke(
		RegisterRoutes,
	),
	healthcheck.Module,
	openapi.Module,
)

func RegisterRoutes(routes *httpservice.Router) {
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.CorrelationIdMiddleware())
	routes.Use(middlewares.CorsMiddleware())

	routes.
		Route("GET /", func(ctx *httpservice.Context) httpservice.Result {
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
