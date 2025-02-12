package handlers

import (
	"github.com/Alias1177/merch-store/internal/handlers/contract"
	contract2 "github.com/Alias1177/merch-store/internal/usecase/contract"
)

type Handler struct {
	userU   contract.UserUsecase
	usecase contract2.BuyUsecase
}

func New(userU contract.UserUsecase, usecase contract2.BuyUsecase) *Handler {
	return &Handler{
		userU:   userU,
		usecase: usecase,
	}
}
