package handlers

import (
	"github.com/Alias1177/merch-store/internal/usecase/contract"
)

type Handler struct {
	userUsecase contract.UserUsecase
	buyUsecase  contract.BuyUsecase
	infoUsecase contract.InfoUsecase
	sendUsecase contract.CoinsUsecase
}

func New(userU contract.UserUsecase, buyUsecase contract.BuyUsecase, infoUsecase contract.InfoUsecase, sendUsecase contract.CoinsUsecase) *Handler {
	return &Handler{
		userUsecase: userU,
		buyUsecase:  buyUsecase,
		infoUsecase: infoUsecase,
		sendUsecase: sendUsecase,
	}
}
