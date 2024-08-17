package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"strings"
	"time"
)

func AuthMiddleware() httpfx.Handler {
	return func(ctx *httpfx.Context) httpfx.Response {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			return httpfx.Response{
				StatusCode: 401,
				Body:       []byte("Authorization header is missing"),
			}
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return httpfx.Response{
				StatusCode: 401,
				Body:       []byte("Token is missing"),
			}
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return httpfx.Response{StatusCode: 401, Body: nil}, nil
			}
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			return httpfx.Response{StatusCode: 401, Body: []byte(err.Error())}
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return httpfx.Response{StatusCode: 401, Body: []byte("Invalid token")}
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return httpfx.Response{StatusCode: 401, Body: []byte("Token is expired")}
			}
		}

		return ctx.Next()
	}
}
