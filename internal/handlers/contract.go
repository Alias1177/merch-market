package handlers

import (
	"context"
	"github.com/Alias1177/merch-store/internal/models"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, reqData models.RegisterRequest) (string, error)
}
