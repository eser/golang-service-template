package middlewares

import (
	"github.com/eser/go-service/pkg/bliss/httpservice"
)

const AccessControlAllowOriginHeader = "Access-Control-Allow-Origin"

func CorsMiddleware() httpservice.Handler {
	return func(ctx *httpservice.Context) httpservice.Result {
		result := ctx.Next()

		ctx.ResponseWriter.Header().Set(AccessControlAllowOriginHeader, "*")

		return result
	}
}
