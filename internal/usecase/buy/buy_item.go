package buy

import (
	"context"
	"fmt"
	"github.com/Alias1177/merch-store/internal/handlers/contract"
)

// BuyUsecaseImpl реализует интерфейс BuyUsecase
type BuyUsecaseImpl struct {
	repo contract.BuyRepo
}

func NewBuyUsecase(repo contract.BuyRepo) *BuyUsecaseImpl {
	return &BuyUsecaseImpl{repo: repo}
}

// Метод покупки предмета
func (u *BuyUsecaseImpl) BuyItem(ctx context.Context, userID int, itemID int) error {
	if err := u.repo.BuyItem(ctx, userID, itemID); err != nil {
		return fmt.Errorf("error processing purchase: %w", err)
	}
	return nil
}
