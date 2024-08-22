package middlewares

import "github.com/eser/go-service/pkg/bliss/httpfx"

func ErrorHandlerMiddleware() httpfx.Handler {
	return func(ctx *httpfx.Context) httpfx.ResponseResult {
		result := ctx.Next()

		return result
	}
}
