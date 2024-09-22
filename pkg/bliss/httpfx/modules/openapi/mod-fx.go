package openapi

import (
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"go.uber.org/fx"
)

type ApiIdentity struct {
	name    string
	version string
}

var FxModule = fx.Module( //nolint:gochecknoglobals
	"openapi",
	fx.Invoke(
		RegisterRoutes,
	),
)

func RegisterRoutes(routes httpfx.Router) {
	routes.
		Route("GET /openapi.json", func(ctx *httpfx.Context) httpfx.Result {
			spec := &ApiIdentity{
				name:    "golang-service",
				version: "0.0.0",
			}

			result := GenerateOpenApiSpec(spec, routes)

			return ctx.Results.Json(result)
		}).
		HasSummary("OpenAPI Spec").
		HasDescription("OpenAPI Spec Endpoint")
}
