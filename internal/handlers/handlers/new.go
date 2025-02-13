package handlers

import (
	"github.com/Alias1177/merch-store/internal/usecase/contract"
)

type Handler struct {
	userUsecase contract.UserUsecase
	buyUsecase  contract.BuyUsecase
	infoUsecase contract.InfoUsecase
}

func New(userU contract.UserUsecase, buyUsecase contract.BuyUsecase, infoUsecase contract.InfoUsecase) *Handler {
	return &Handler{
		userUsecase: userU,
		buyUsecase:  buyUsecase,
		infoUsecase: infoUsecase,
	}
}
