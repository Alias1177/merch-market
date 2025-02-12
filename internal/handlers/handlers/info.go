// handlers/info.go
package handlers

import (
	"encoding/json"
	"github.com/Alias1177/merch-store/internal/constants"
	"github.com/Alias1177/merch-store/internal/models"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func InfoHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Context values: %+v", r.Context())
		contextValue := r.Context().Value(constants.UserIDContextKey)
		log.Printf("UserID from context: %+v", contextValue)

		userID, ok := contextValue.(int)
		if !ok {
			log.Printf("Type assertion failed. Type of contextValue: %T", contextValue)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Начинаем транзакцию
		tx, err := db.Beginx()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// Получаем баланс монет
		var coins int
		err = tx.Get(&coins, "SELECT coins FROM users WHERE id = $1", userID)
		if err != nil {
			http.Error(w, "Failed to get user coins", http.StatusInternalServerError)
			return
		}

		// Получаем инвентарь
		var inventory []models.InventoryItem
		err = tx.Select(&inventory, `
    SELECT i.name, inv.quantity 
    FROM inventory inv
    JOIN items i ON inv.item_id = i.id 
    WHERE inv.user_id = $1`, userID)
		if err != nil {
			http.Error(w, "Failed to get inventory", http.StatusInternalServerError)
			return
		}

		// Получаем историю полученных монет
		var received []models.ReceivedTransaction
		err = tx.Select(&received, `
    SELECT u.username, t.amount 
    FROM transactions t 
    JOIN users u ON t.sender_id = u.id 
    WHERE t.receiver_id = $1`, userID)
		if err != nil {
			http.Error(w, "Failed to get received history", http.StatusInternalServerError)
			return
		}

		// Получаем историю отправленных монет
		var sent []models.SentTransaction
		err = tx.Select(&sent, `
    SELECT u.username, t.amount 
    FROM transactions t 
    JOIN users u ON t.receiver_id = u.id 
    WHERE t.sender_id = $1`, userID)
		if err != nil {
			http.Error(w, "Failed to get sent history", http.StatusInternalServerError)
			return
		}

		// Формируем ответ
		response := models.InfoResponse{
			Coins:     coins,
			Inventory: inventory,
			CoinHistory: models.CoinHistoryDetails{
				Received: received,
				Sent:     sent,
			},
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")

		// Отправляем ответ
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
