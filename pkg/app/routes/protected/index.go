package protected

import (
	"fmt"
	"net/http"

	"github.com/eser/go-service/pkg/app/config"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/httpfx/middlewares"
)

func IndexRoutes(routes httpfx.Router, appConfig *config.AppConfig) {
	routes.
		Route("GET /protected", middlewares.AuthMiddleware(), func(ctx *httpfx.Context) httpfx.Result {
			message := fmt.Sprintf("Hello from %s! this endpoint is protected!", appConfig.AppName)

			return ctx.Results.PlainText(message)
		}).
		HasSummary("Protected page").
		HasDescription("A page protected with JWT auth.").
		HasResponse(http.StatusOK)
}
