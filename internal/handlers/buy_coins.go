package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/Alias1177/merch-store/pkg/logger"
	"net/http"
	"strconv"

	"github.com/Alias1177/merch-store/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type APIErrorResponse struct {
	Errors string `json:"errors"`
}

func BuyHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.ColorLogger()
		// Получаем item из пути /api/buy/{item}
		itemIDStr := chi.URLParam(r, "item")
		itemID, err := strconv.Atoi(itemIDStr)
		if err != nil || itemID <= 0 {
			// Если ID предмета некорректный, возвращаем ошибку
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIErrorResponse{Errors: "Invalid item ID"})
			return
		}

		// Получаем userID из контекста
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			// Если пользователь не авторизован, возвращаем ошибку
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(APIErrorResponse{Errors: "Unauthorized"})
			return
		}

		// Открываем транзакцию
		tx, err := db.Beginx()
		if err != nil {
			// Если не удалось начать транзакцию, возвращаем ошибку
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIErrorResponse{Errors: "Failed to start transaction"})
			return
		}

		// Проверяем, существует ли предмет и получаем его цену
		var price int
		err = tx.Get(&price, "SELECT price FROM items WHERE id = $1", itemID)
		if err == sql.ErrNoRows {
			// Если предмет не найден, откатываем транзакцию и возвращаем ошибку
			tx.Rollback()
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(APIErrorResponse{Errors: "Item not found"})
			return
		} else if err != nil {
			// В случае общей ошибки при выполнении запроса возвращаем ошибку
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIErrorResponse{Errors: "Failed to query item"})
			return
		}

		// Получаем баланс пользователя
		var coins int
		err = tx.Get(&coins, "SELECT coins FROM users WHERE id = $1", userID)
		if err == sql.ErrNoRows {
			// Если пользователь не найден, откатываем транзакцию и возвращаем ошибку
			tx.Rollback()
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(APIErrorResponse{Errors: "User not found"})
			return
		} else if err != nil {
			// В случае ошибки при запросе баланса пользователя возвращаем ошибку
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIErrorResponse{Errors: "Failed to query user balance"})
			return
		}

		// Проверяем, хватает ли монет у пользователя для покупки предмета
		if coins < price {
			// Если монет недостаточно, откатываем транзакцию и возвращаем ошибку
			tx.Rollback()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIErrorResponse{Errors: "Not enough coins"})
			return
		}

		// Списываем стоимость предмета с баланса пользователя
		_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id = $2", price, userID)
		if err != nil {
			// В случае ошибки при обновлении баланса откатываем транзакцию и возвращаем ошибку
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIErrorResponse{Errors: "Failed to update user balance"})
			return
		}

		// Обновляем инвентарь пользователя или добавляем новый предмет
		_, err = tx.Exec(`
			INSERT INTO inventory (user_id, item_id, quantity)
			VALUES ($1, $2, 1)
			ON CONFLICT (user_id, item_id)
			DO UPDATE SET quantity = inventory.quantity + 1`,
			userID, itemID)
		if err != nil {
			// В случае ошибки при обновлении инвентаря откатываем транзакцию и возвращаем ошибку
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIErrorResponse{Errors: "Failed to update inventory"})
			return
		}

		// Завершаем транзакцию
		if err := tx.Commit(); err != nil {
			// В случае ошибки при завершении транзакции возвращаем ошибку
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(APIErrorResponse{Errors: "Failed to commit transaction"})
			return
		}

		// Возвращаем успешный ответ
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Item purchased successfully",
			"itemID":  itemID,
		})
	}
}
