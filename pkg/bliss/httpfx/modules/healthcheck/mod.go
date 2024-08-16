package healthcheck

import (
	"net/http"

	"github.com/eser/go-service/pkg/bliss/httpfx"
	"go.uber.org/fx"
)

var Module = fx.Module( //nolint:gochecknoglobals
	"healthcheck",
	fx.Invoke(
		RegisterRoutes,
	),
)

func RegisterRoutes(routes httpfx.Router) {
	routes.
		Route("GET /health-check", func(ctx *httpfx.Context) httpfx.Response {
			return ctx.Results.Ok()
		}).
		HasSummary("Health Check").
		HasDescription("Health Check Endpoint").
		HasResponse(http.StatusNoContent)
}