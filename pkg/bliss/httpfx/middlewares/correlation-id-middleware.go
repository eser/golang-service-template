package middlewares

import (
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/lib"
)

const CorrelationIdHeader = "X-Correlation-Id"

func CorrelationIdMiddleware() httpfx.Handler {
	return func(ctx *httpfx.Context) httpfx.Response {
		correlationId := ctx.Request.Header.Get(CorrelationIdHeader)
		if correlationId == "" {
			correlationId = lib.GenerateUniqueId()
		}

		result := ctx.Next()

		ctx.ResponseWriter.Header().Set(CorrelationIdHeader, correlationId)

		return result
	}
}
