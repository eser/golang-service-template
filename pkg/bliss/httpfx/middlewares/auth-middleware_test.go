package middlewares

import (
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func createToken(secret string, exp time.Time) string {
	claims := jwt.MapClaims{
		"exp": exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		expectedCode int
	}{
		{
			name:         "No Authorization Header",
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Invalid Token Format",
			token:        "InvalidToken",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Expired Token",
			token:        createToken("secret", time.Now().Add(-time.Hour)),
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Valid Token",
			token:        createToken("secret", time.Now().Add(time.Hour)),
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			res := httptest.NewRecorder()
			ctx := httpfx.Context{
				Request:        req,
				ResponseWriter: res,
			}

			middleware := AuthMiddleware()
			result := middleware(&ctx)

			if result.StatusCode != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, result.StatusCode)
			}

			if tt.expectedCode == http.StatusOK {
				if res.Header().Get("Authorization") == "" {
					t.Error("Authorization header is missing")
				}
			}
		})
	}
}
