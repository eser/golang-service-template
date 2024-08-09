package healthcheck

import (
	"net/http"

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
			return ctx.Results.Ok()
		}).
		HasSummary("Health Check").
		HasDescription("Health Check Endpoint").
		HasResponse(http.StatusNoContent)
}
