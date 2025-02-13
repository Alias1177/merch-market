package contract

import (
	"context"

	"github.com/Alias1177/merch-store/internal/models"
)

type DBRepo interface {
	CreateUser(ctx context.Context, username, passwordHash string, coins int) (*models.User, error)
}
type UserUsecase interface {
	CreateUser(ctx context.Context, reqData models.RegisterRequest) (string, error)
}
type BuyRepo interface {
	BuyItem(ctx context.Context, userID, itemID int) error
}
type BuyUsecase interface {
	BuyItem(ctx context.Context, userID, itemID int) error
}
type InfoUsecase interface {
	GetUserInfo(ctx context.Context, userID int) (*models.InfoResponse, error)
}
type InfoRepository interface {
	GetUserInfo(ctx context.Context, userID int) (*models.InfoResponse, error)
}
type CoinsRepository interface {
	SendCoins(ctx context.Context, senderID int, receiverUsername string, amount int) error
}
type CoinsUsecase interface {
	SendCoins(ctx context.Context, senderID int, receiverUsername string, amount int) error
}
