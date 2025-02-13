package handlers

import (
	"encoding/json"
	"github.com/Alias1177/merch-store/internal/middleware"
	"github.com/Alias1177/merch-store/internal/models"
	"net/http"
)

func (h *Handler) HandleSendCoins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	senderID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	err = h.sendUsecase.SendCoins(r.Context(), senderID, req.ToUser, req.Amount)
	if err != nil {
		http.Error(w, "Failed to send coins: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Coins sent successfully"})
}
