package middlewares

import (
	"log"
	"time"

	"github.com/eser/go-service/pkg/bliss/httpfx"
)

const ResponseTimeHeader = "X-Request-Time"

func ResponseTimeMiddleware() httpfx.Handler {
	return func(ctx *httpfx.Context) httpfx.Result {
		startTime := time.Now()

		result := ctx.Next()

		duration := time.Since(startTime)
		// durationText := strconv.FormatInt(duration.Milliseconds(), 10) + "ms"
		durationText := duration.String()

		log.Println(ctx.Request.Method, ctx.Request.URL.Path, duration)

		ctx.ResponseWriter.Header().Set(ResponseTimeHeader, durationText)

		return result
	}
}
