package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Alias1177/merch-store/internal/constants"
)

func TestJWTMiddleware(t *testing.T) {
	secretKey := "supersecretkey"
	handler := func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(constants.UserIDContextKey)
		if userID == nil {
			http.Error(w, "UserID not set in context", http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Missing Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Authorization header format",
			authHeader:     "InvalidFormat",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Signing Method",
			authHeader:     createTokenWithInvalidMethod(secretKey),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid or expired token",
			authHeader:     "Bearer invalid.token.value",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid token with missing user_id claim",
			authHeader:     createToken(secretKey, nil),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.authHeader != "" {
				req.Header.Set("Authorization", test.authHeader)
			}
			rr := httptest.NewRecorder()

			middleware := JWTMiddleware(secretKey)
			middleware(http.HandlerFunc(handler)).ServeHTTP(rr, req)

			if rr.Code != test.expectedStatus {
				t.Errorf("expected status %d, got %d", test.expectedStatus, rr.Code)
			}
		})
	}
}

// Вспомогательная функция для создания токена с корректным или отсутствующим user_id
func createToken(secretKey string, claims map[string]interface{}) string {
	tokenClaims := jwt.MapClaims{}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}
	return "Bearer " + signedToken
}

// Вспомогательная функция для создания токена с неподдерживаемым методом подписи
func createTokenWithInvalidMethod(secretKey string) string {
	// Генерация фиктивного RSA закрытого ключа для SigningMethodRS256
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate RSA key: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id": 123,
	})

	// Подписываем токен с использованием закрытого RSA ключа
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	return "Bearer " + signedToken
}
