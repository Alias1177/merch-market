package handlers

import (
	"encoding/json"
	"github.com/Alias1177/merch-store/internal/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) HandleBuy(w http.ResponseWriter, r *http.Request) {
	// Извлечение itemID из параметров URL
	itemIDStr := chi.URLParam(r, "item")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil || itemID <= 0 {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Получение userID из контекста
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Выполнение бизнес-логики покупки
	if err := h.buyUsecase.BuyItem(r.Context(), userID, itemID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Item purchased successfully!"})
}
