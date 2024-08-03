package healthcheck

import (
	"github.com/eser/go-service/pkg/bliss/httpservice"
	"go.uber.org/fx"
)

var Module = fx.Module( //nolint:gochecknoglobals
	"healthcheck",
	fx.Invoke(
		RegisterRoutes,
	),
)

func RegisterRoutes(routes *httpservice.Router) {
	routes.
		Route("GET /health-check", func(ctx *httpservice.Context) httpservice.Result {
			return ctx.Results.Ok().WithBody("OK")
		}).
		WithSummary("Health Check").
		WithDescription("Health Check Endpoint")
}
