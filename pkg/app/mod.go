package app

import (
	"github.com/eser/go-service/pkg/bliss"
	"github.com/eser/go-service/pkg/bliss/httpservice"
	"github.com/eser/go-service/pkg/bliss/httpservice/middlewares"
	"github.com/eser/go-service/pkg/bliss/httpservice/modules/healthcheck"
	"go.uber.org/fx"
)

var Module = fx.Module( //nolint:gochecknoglobals
	"app",
	fx.Invoke(
		RegisterRoutes,
	),
	healthcheck.Module,
)

func RegisterRoutes(routes *httpservice.Router) {
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.CorrelationIdMiddleware())

	routes.
		Route("GET /", func(ctx *httpservice.Context) httpservice.Result {
			return ctx.Results.PlainText("Hello, World!")
		}).
		WithSummary("Homepage").
		WithDescription("This is the homepage of the service.")
}

func New() *fx.App {
	return fx.New(
		// fx.WithLogger(bliss.GetFxLogger),
		bliss.Module,
		Module,
	)
}
