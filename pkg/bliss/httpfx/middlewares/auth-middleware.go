package middlewares

import (
	"context"
	"strings"
	"time"

	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/results"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthClaims httpfx.ContextKey = "claims"
)

var ErrInvalidSigningMethod = results.Define("ERRBHMA001", "Invalid signing method")

func AuthMiddleware() httpfx.Handler {
	return func(ctx *httpfx.Context) httpfx.Result {
		tokenString, hasToken := getBearerToken(ctx)

		if !hasToken {
			return ctx.Results.Unauthorized("No suitable authorization header found")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidSigningMethod.New()
			}

			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			return ctx.Results.Unauthorized(err.Error())
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return ctx.Results.Unauthorized("Invalid token")
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return ctx.Results.Unauthorized("Token is expired")
			}
		}

		ctx.UpdateContext(context.WithValue(
			ctx.Request.Context(),
			AuthClaims,
			claims,
		))

		return ctx.Next()
	}
}

func getBearerToken(ctx *httpfx.Context) (string, bool) {
	for _, authHeader := range ctx.Request.Header["Authorization"] {
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			return tokenString, true
		}
	}

	return "", false
}
