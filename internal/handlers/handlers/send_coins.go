package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Alias1177/merch-store/internal/middleware"
	"github.com/Alias1177/merch-store/internal/models"
)

func (h *Handler) HandleSendCoins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	senderID, err := middleware.GetUserID(r.Context())
	if err != nil {
		slog.Error("Unauthorized", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Invalid request format", "error", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if req.Amount <= 0 {
		slog.Error("Invalid amount", "amount", req.Amount)
		http.Error(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	if req.ToUser == "" {
		slog.Error("Empty receiver username")
		http.Error(w, "Receiver username cannot be empty", http.StatusBadRequest)
		return
	}

	err = h.sendUsecase.SendCoins(r.Context(), senderID, req.ToUser, req.Amount)
	if err != nil {
		slog.Error("Failed to send coins", "error", err)

		switch err.Error() {
		case "user not found":
			http.Error(w, "Receiver not found", http.StatusNotFound)
		case "not enough coins":
			http.Error(w, "Insufficient funds", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Coins sent successfully",
	}); err != nil {
		slog.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
