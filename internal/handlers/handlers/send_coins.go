package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Alias1177/merch-store/internal/middleware"
	"github.com/Alias1177/merch-store/internal/models"
	"net/http"
)

func (h *Handler) HandleSendCoins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	senderID, err := middleware.GetUserID(r.Context())
	if err != nil {
		fmt.Errorf("Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Errorf("Invalid request format", http.StatusBadRequest)
		return
	}

	err = h.sendUsecase.SendCoins(r.Context(), senderID, req.ToUser, req.Amount)
	if err != nil {
		fmt.Errorf("Failed to send coins", http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Coins sent successfully"})
}
