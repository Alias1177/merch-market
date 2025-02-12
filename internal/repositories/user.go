package repositories

import (
	"errors"
	"github.com/Alias1177/merch-store/pkg/logger"
	"github.com/jmoiron/sqlx"
)

// Ошибка конфликта уникальности
var ErrUserAlreadyExists = errors.New("user already exists")

// Упрощённая структура User
type User struct {
	ID           int    `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	Coins        int    `db:"coins"`
}

// CreateUser добавляет нового пользователя в базу данных.
func CreateUser(db *sqlx.DB, username, passwordHash string, coins int) (*User, error) {
	logger.ColorLogger()
	// Упрощённый SQL для добавления нового пользователя
	query := `
		INSERT INTO users (username, password_hash, coins)
		VALUES ($1, $2, $3)
		RETURNING id, username, password_hash, coins
	`

	// Объект для сохранения результата
	user := &User{}

	// Выполнение запроса и возврат результата
	err := db.QueryRow(query, username, passwordHash, coins).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Coins,
	)
	if err != nil {
		// Проверяем, если ошибка вызвана нарушением уникальности
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return user, nil
}
