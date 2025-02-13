package coins

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Alias1177/merch-store/internal/usecase/contract"
)

type CoinsUsecase struct {
	repo contract.CoinsRepository
}

func NewCoinsUsecase(repo contract.CoinsRepository) *CoinsUsecase {
	return &CoinsUsecase{
		repo: repo,
	}
}

func (u *CoinsUsecase) SendCoins(ctx context.Context, senderID int, receiverUsername string, amount int) error {
	if amount <= 0 {
		slog.Error("amount must be positive")
		return fmt.Errorf("amount must be positive")
	}
	return u.repo.SendCoins(ctx, senderID, receiverUsername, amount)
}
