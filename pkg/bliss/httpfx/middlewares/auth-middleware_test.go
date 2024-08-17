package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/eser/go-service/pkg/bliss/httpfx"
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
	t.Parallel()

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
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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

			// FIXME(@eser): temporarily disabled due to understanding the expected behavior
			// if tt.expectedCode == http.StatusOK || tt.expectedCode == http.StatusNoContent {
			// 	if res.Header().Get(" ") == "" {
			// 		t.Error("Authorization header is missing")
			// 	}
			// }
		})
	}
}
