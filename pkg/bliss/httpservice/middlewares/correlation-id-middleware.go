package middlewares

import (
	"github.com/eser/go-service/pkg/bliss/httpservice"
)

func CorrelationIdMiddleware() httpservice.Handler {
	return func(ctx *httpservice.Context) httpservice.Result {
		correlationId := ctx.Request.Header.Get("X-Correlation-Id")

		result := ctx.Next()

		if correlationId != "" {
			ctx.ResponseWriter.Header().Set("X-Correlation-Id", correlationId)
		}

		return result
	}
}
