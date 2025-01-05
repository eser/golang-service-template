package healthcheck

import (
	"net/http"

	"github.com/eser/go-service/pkg/bliss/httpfx"
)

func RegisterHttpRoutes(routes httpfx.Router, config *httpfx.Config) {
	if !config.HealthCheckEnabled {
		return
	}

	routes.
		Route("GET /health-check", func(ctx *httpfx.Context) httpfx.Result {
			return ctx.Results.Ok()
		}).
		HasSummary("Health Check").
		HasDescription("Health Check Endpoint").
		HasResponse(http.StatusNoContent)
}
