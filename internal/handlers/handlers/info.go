package handlers

import (
	"encoding/json"
	"github.com/Alias1177/merch-store/internal/middleware"
	"net/http"
)

func (h *Handler) HandleInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	info, err := h.infoUsecase.GetUserInfo(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
