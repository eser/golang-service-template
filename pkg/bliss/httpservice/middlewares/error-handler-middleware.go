package middlewares

import "github.com/eser/go-service/pkg/bliss/httpservice"

func ErrorHandlerMiddleware() httpservice.Handler {
	return func(ctx *httpservice.Context) httpservice.Result {
		result := ctx.Next()

		return result
	}
}
