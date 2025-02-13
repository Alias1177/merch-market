package info

import (
	"context"
	"github.com/Alias1177/merch-store/internal/models"
	"github.com/Alias1177/merch-store/internal/usecase/contract"
)

type InfoUsecase struct {
	repo contract.InfoRepository
}

func NewInfoUsecase(repo contract.InfoRepository) *InfoUsecase {
	return &InfoUsecase{
		repo: repo,
	}
}

func (u *InfoUsecase) GetUserInfo(ctx context.Context, userID int) (*models.InfoResponse, error) {
	return u.repo.GetUserInfo(ctx, userID)
}
