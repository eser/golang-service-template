package middlewares

import (
	"log"
	"time"

	"github.com/eser/go-service/pkg/bliss/httpservice"
)

func ResponseTimeMiddleware() httpservice.Handler {
	return func(ctx *httpservice.Context) httpservice.Result {
		startTime := time.Now()

		result := ctx.Next()

		duration := time.Since(startTime)
		// durationText := strconv.FormatInt(duration.Milliseconds(), 10) + "ms"
		durationText := duration.String()

		log.Println(ctx.Request.Method, ctx.Request.URL.Path, duration)

		ctx.ResponseWriter.Header().Set("X-Request-Time", durationText)

		return result
	}
}
