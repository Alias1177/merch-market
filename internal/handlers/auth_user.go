package handlers

import (
	"encoding/json"
	"github.com/Alias1177/merch-store/config/config"
	"github.com/Alias1177/merch-store/pkg/logger"
	"net/http"

	"github.com/Alias1177/merch-store/internal/middleware/jwt"
	"github.com/Alias1177/merch-store/internal/repositories"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest представляет тело запроса для регистрации пользователя
type RegisterRequest struct {
	Username string `json:"username"` // Имя пользователя
	Password string `json:"password"` // Пароль
}

// TokenResponse представляет ответ, содержащий JWT токен
type TokenResponse struct {
	Token string `json:"token"` // JWT токен
}

// Загружаем конфигурацию из файла
var cfg = config.Load("./config/config.yaml")

// RegisterHandler обрабатывает запросы регистрации пользователя
func RegisterHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.ColorLogger()
		// Декодируем тело запроса в структуру RegisterRequest
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest) // Ошибка формата запроса
			return
		}

		// Хэшируем пароль с помощью bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError) // Ошибка хэширования пароля
			return
		}

		// Создаём пользователя в базе данных
		user, err := repositories.CreateUser(db, req.Username, string(hashedPassword), 1000) // 1000 начальных монет
		if err != nil {
			// Проверяем, если пользователь уже существует
			if err == repositories.ErrUserAlreadyExists {
				http.Error(w, "User already exists", http.StatusConflict) // Конфликт - пользователь уже существует
				return
			}
			http.Error(w, "Failed to create user", http.StatusInternalServerError) // Ошибка создания пользователя
			return
		}

		// Генерируем JWT токен для нового пользователя
		token, err := middleware.GenerateJWT(user.ID, user.Username, cfg.JWT.Secret)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError) // Ошибка генерации токена
			return
		}

		// Возвращаем токен в ответе клиенту
		json.NewEncoder(w).Encode(TokenResponse{Token: token})
	}
}
