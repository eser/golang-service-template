package middlewares

import (
	"github.com/eser/go-service/pkg/bliss/httpservice"
	"github.com/eser/go-service/pkg/bliss/lib"
)

const CorrelationIdHeader = "X-Correlation-Id"

func CorrelationIdMiddleware() httpservice.Handler {
	return func(ctx *httpservice.Context) httpservice.Result {
		correlationId := ctx.Request.Header.Get(CorrelationIdHeader)
		if correlationId == "" {
			correlationId = lib.GenerateUniqueId()
		}

		result := ctx.Next()

		ctx.ResponseWriter.Header().Set(CorrelationIdHeader, correlationId)

		return result
	}
}
