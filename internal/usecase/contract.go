package usecase

import (
	"context"
	"github.com/Alias1177/merch-store/internal/models"
)

type DBRepo interface {
	CreateUser(ctx context.Context, username, passwordHash string, coins int) (*models.User, error)
}
