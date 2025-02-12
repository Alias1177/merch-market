package models

// Упрощённая структура User
type User struct {
	ID           int    `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	Coins        int    `db:"coins"`
}

// RegisterRequest представляет тело запроса для регистрации пользователя

// TokenResponse представляет ответ, содержащий JWT токен
type TokenResponse struct {
	Token string `json:"token"` // JWT токен
}
