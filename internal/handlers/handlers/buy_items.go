package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Alias1177/merch-store/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) HandleBuy(w http.ResponseWriter, r *http.Request) {
	// Извлечение itemID из параметров URL
	itemIDStr := chi.URLParam(r, "item")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil || itemID <= 0 {
		slog.Error("Invalid item ID")
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Получение userID из контекста
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		slog.Error("Unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Выполнение бизнес-логики покупки
	if err := h.buyUsecase.BuyItem(r.Context(), userID, itemID); err != nil {
		slog.Error("Failed to buy item: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Item purchased successfully!"}); err != nil {
		slog.Error("Failed to encode response: " + err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
