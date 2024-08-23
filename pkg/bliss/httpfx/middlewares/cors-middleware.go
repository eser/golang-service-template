package middlewares

import (
	"github.com/eser/go-service/pkg/bliss/httpfx"
)

const AccessControlAllowOriginHeader = "Access-Control-Allow-Origin"

func CorsMiddleware() httpfx.Handler {
	return func(ctx *httpfx.Context) httpfx.Result {
		result := ctx.Next()

		ctx.ResponseWriter.Header().Set(AccessControlAllowOriginHeader, "*")

		return result
	}
}
