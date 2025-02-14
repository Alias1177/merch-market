package middleware

import (
	"context"
	"errors"

	"github.com/Alias1177/merch-store/internal/constants"
)

func GetUserID(ctx context.Context) (int, error) {
	userID := ctx.Value(constants.UserIDContextKey)
	if userID == nil {
		return 0, errors.New("userID not found in context")
	}

	id, ok := userID.(int)
	if !ok {
		return 0, errors.New("invalid userID format")
	}

	return id, nil
}
