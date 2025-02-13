package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Alias1177/merch-store/internal/models"
)

// RegisterHandler обрабатывает запросы регистрации пользователя
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Invalid request format")
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	token, err := h.userUsecase.CreateUser(r.Context(), req)
	if err != nil {
		slog.Error("Failed to create user:")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(models.TokenResponse{Token: token}); err != nil {
		slog.Error("Error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
