package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	middleware "github.com/Alias1177/merch-store/internal/middleware/jwt"
	"github.com/Alias1177/merch-store/internal/models"
	"github.com/Alias1177/merch-store/internal/usecase/contract"
	"github.com/Alias1177/merch-store/pkg"
)

type UserUsecase struct {
	dbR    contract.DBRepo
	secret string
}

func New(dbR contract.DBRepo, secret string) *UserUsecase {
	return &UserUsecase{
		dbR:    dbR,
		secret: secret,
	}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, reqData models.RegisterRequest) (string, error) {
	// Хэшируем пароль с помощью bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("error hashing password:")
		return "", fmt.Errorf("error hashing password: %v", err)
	}

	// Создаём пользователя в базе данных
	user, err := uc.dbR.CreateUser(ctx, reqData.Username, string(hashedPassword), 1000) // 1000 начальных монет
	if err != nil {
		// Проверяем, если пользователь уже существует
		if errors.Is(err, pkg.ErrUserAlreadyExists) {
			slog.Error("user already exists:")
			return "", errors.New("user already exists")
		}
		slog.Error("error creating user:")
		return "", fmt.Errorf("error creating user: %v", err)
	}

	// Генерируем JWT токен для нового пользователя
	token, err := middleware.GenerateJWT(user.ID, user.Username, uc.secret)
	if err != nil {
		slog.Error("error generating JWT token:")
		return "", fmt.Errorf("error generating JWT token: %v", err)
	}
	return token, nil
}
