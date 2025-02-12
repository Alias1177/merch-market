package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/Alias1177/merch-store/internal/middleware"
	"github.com/Alias1177/merch-store/pkg/logger"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type APIResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func SendCoinsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.ColorLogger()
		w.Header().Set("Content-Type", "application/json")

		// Получаем ID отправителя из контекста
		senderID, err := middleware.GetUserID(r.Context())
		if err != nil {
			json.NewEncoder(w).Encode(APIResponse{Error: "Unauthorized"})
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Декодируем запрос
		var req SendCoinRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIResponse{Error: "Invalid request format"})
			return
		}

		// Валидация запроса
		if req.Amount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIResponse{Error: "Amount must be positive"})
			return
		}

		// Начинаем транзакцию
		tx, err := db.Beginx()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{Error: "Internal server error"})
			return
		}
		defer tx.Rollback()

		// Получаем ID получателя
		var receiverID int
		err = tx.Get(&receiverID, "SELECT id FROM users WHERE username = $1", req.ToUser)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIResponse{Error: "Recipient not found"})
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{Error: "Internal server error"})
			return
		}

		// Проверяем баланс отправителя
		var senderCoins int
		err = tx.Get(&senderCoins, "SELECT coins FROM users WHERE id = $1", senderID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{Error: "Internal server error"})
			return
		}

		if senderCoins < req.Amount {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIResponse{Error: "Insufficient coins"})
			return
		}

		// Выполняем перевод
		_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id = $2", req.Amount, senderID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{Error: "Failed to update sender balance"})
			return
		}

		_, err = tx.Exec("UPDATE users SET coins = coins + $1 WHERE id = $2", req.Amount, receiverID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{Error: "Failed to update recipient balance"})
			return
		}

		// Записываем транзакцию
		_, err = tx.Exec(`
			INSERT INTO transactions (sender_id, receiver_id, amount)
			VALUES ($1, $2, $3)
		`, senderID, receiverID, req.Amount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{Error: "Failed to record transaction"})
			return
		}

		// Подтверждаем транзакцию
		if err := tx.Commit(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIResponse{Error: "Failed to complete transaction"})
			return
		}

		// Отправляем успешный ответ
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(APIResponse{Message: "Coins sent successfully"})
	}
}
