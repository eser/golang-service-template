package testhttp

import (
	"fmt"
	"net/http"

	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/httpfx/middlewares"
)

func RegisterHttpRoutes(routes httpfx.Router, appConfig *AppConfig) {
	routes.
		Route("GET /protected", middlewares.AuthMiddleware(), func(ctx *httpfx.Context) httpfx.Result {
			v := []string{
				"1234",
				"eser.live",
			}

			return ctx.Results.Json(v)
		}).
		HasSummary("Protected page").
		HasDescription("A page protected with JWT auth.").
		HasResponse(http.StatusOK)

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
