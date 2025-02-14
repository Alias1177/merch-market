package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Alias1177/merch-store/internal/constants"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware возвращает middleware для проверки JWT токенов
func JWTMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Извлекаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			// Проверяем формат токена (Bearer <token>)
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Разбираем и проверяем токен с использованием секретного ключа
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Убеждаемся, что используется метод подписи HMAC SHA256
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(secretKey), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Извлекаем user_id из claims токена
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["user_id"].(float64) // В JWT числа возвращаются как float64
			if !ok {
				http.Error(w, "Invalid token payload", http.StatusUnauthorized)
				return
			}

			// Сохраняем userID в контексте для использования в хендлерах
			// В jwtParse.go
			log.Printf("Setting userID in context: %v", int(userID))
			ctx := context.WithValue(r.Context(), constants.UserIDContextKey, int(userID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
