package repositories

import (
	"context"

	"github.com/Alias1177/merch-store/internal/models"
	"github.com/Alias1177/merch-store/pkg"
	"github.com/Alias1177/merch-store/pkg/logger"
)

func (r *Repository) CreateUser(ctx context.Context, username, passwordHash string, coins int) (*models.User, error) {
	logger.ColorLogger()
	// Упрощённый SQL для добавления нового пользователя
	query := `
		INSERT INTO users (username, password_hash, coins)
		VALUES ($1, $2, $3)
		RETURNING id, username, password_hash, coins
	`

	// Объект для сохранения результата
	user := &models.User{}

	// Выполнение запроса и возврат результата
	err := r.conn.QueryRow(query, username, passwordHash, coins).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Coins,
	)
	if err != nil {
		// Проверяем, если ошибка вызвана нарушением уникальности
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return nil, pkg.ErrUserAlreadyExists
		}
		return nil, err
	}

	return user, nil
}
