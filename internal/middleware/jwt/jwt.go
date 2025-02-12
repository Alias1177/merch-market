package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Генерация JWT токена
func GenerateJWT(userID int, username, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
