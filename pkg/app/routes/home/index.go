package home

import (
	"fmt"
	"net/http"

	"github.com/eser/go-service/pkg/app/config"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/httpfx/middlewares"
)

func IndexRoutes(routes httpfx.Router, appConfig *config.AppConfig) {
	routes.
		Route("GET /", func(ctx *httpfx.Context) httpfx.Result {
			message := fmt.Sprintf(
				"Hello %s (%s) from %s!",
				ctx.Request.Context().Value(middlewares.ClientAddr),
				ctx.Request.Context().Value(middlewares.ClientAddrOrigin),
				appConfig.AppName,
			)

			return ctx.Results.PlainText(message)
		}).
		HasSummary("Homepage").
		HasDescription("This is the homepage of the service.").
		HasResponse(http.StatusOK)
}
