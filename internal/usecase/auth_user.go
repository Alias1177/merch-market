package usecase

import (
	"context"
	"errors"
	"fmt"
	middleware "github.com/Alias1177/merch-store/internal/middleware/jwt"
	"github.com/Alias1177/merch-store/internal/models"
	"github.com/Alias1177/merch-store/pkg"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	dbR    DBRepo
	secret string
}

func New(dbR DBRepo, secret string) *UserUsecase {
	return &UserUsecase{
		dbR:    dbR,
		secret: secret,
	}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, reqData models.RegisterRequest) (string, error) {
	// Хэшируем пароль с помощью bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}

	// Создаём пользователя в базе данных
	user, err := uc.dbR.CreateUser(ctx, reqData.Username, string(hashedPassword), 1000) // 1000 начальных монет
	if err != nil {
		// Проверяем, если пользователь уже существует
		if errors.Is(err, pkg.ErrUserAlreadyExists) {
			return "", errors.New("user already exists")

		}
		return "", fmt.Errorf("error creating user: %v", err)
	}

	// Генерируем JWT токен для нового пользователя
	token, err := middleware.GenerateJWT(user.ID, user.Username, uc.secret)
	if err != nil {
		return "", fmt.Errorf("error generating JWT token: %v", err)

	}
	return token, nil

}
