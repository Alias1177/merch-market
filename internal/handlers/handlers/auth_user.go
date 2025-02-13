package handlers

import (
	"encoding/json"
	"github.com/Alias1177/merch-store/internal/models"
	"github.com/Alias1177/merch-store/pkg/logger"
	"net/http"
)

// RegisterHandler обрабатывает запросы регистрации пользователя
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	logger.ColorLogger()

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	token, err := h.userUsecase.CreateUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(models.TokenResponse{Token: token}); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
