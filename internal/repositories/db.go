package repositories

import (
	"context"
	"github.com/Alias1177/merch-store/internal/models"
	"github.com/Alias1177/merch-store/pkg/logger"
	"log"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Alias1177/merch-store/pkg"
)

type Repository struct {
	conn *sqlx.DB
}

func New(dsn string) *Repository {
	logger.ColorLogger()
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		slog.Error(pkg.DbError)
		log.Fatalf("Unable to connect to the database: %v", err)
	}
	log.Println("Successfully connected to the database.")
	return &Repository{db}
}

// CreateUser добавляет нового пользователя в базу данных.
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

func (r *Repository) Close() {
	logger.ColorLogger()
	err := r.conn.Close()
	if err != nil {
		slog.Error(pkg.DbError)
		log.Fatalf("Unable to close the database connection: %v", err)
	}
}
