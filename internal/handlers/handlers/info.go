package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Alias1177/merch-store/internal/middleware"
)

func (h *Handler) HandleInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		slog.Error("Unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	info, err := h.infoUsecase.GetUserInfo(r.Context(), userID)
	if err != nil {
		slog.Error("Failed to get user info: " + err.Error())
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
	}
}
