package app

import (
	"github.com/eser/go-service/pkg/app/healthcheck"
	"github.com/eser/go-service/pkg/bliss"
	"github.com/eser/go-service/pkg/bliss/httpservice"
	"github.com/eser/go-service/pkg/bliss/httpservice/middlewares"
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
}

func New() *fx.App {
	return fx.New(
		// fx.WithLogger(bliss.GetFxLogger),
		bliss.Module,
		Module,
	)
}
