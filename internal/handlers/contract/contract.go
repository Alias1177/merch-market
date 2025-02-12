package contract

import (
	"context"
	"github.com/Alias1177/merch-store/internal/models"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, reqData models.RegisterRequest) (string, error)
}

type BuyRepo interface {
	BuyItem(ctx context.Context, userID, itemID int) error
}
