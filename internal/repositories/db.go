package repositories

import (
	"context"
	"log"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Alias1177/merch-store/pkg"
)

type Repository struct {
	conn *sqlx.DB
}

func New(ctx context.Context, dsn string) *Repository {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		slog.Error(pkg.DbError)
		log.Fatalf("Unable to connect to the database: %v", err)
	}
	log.Println("Successfully connected to the database.")
	return &Repository{db}
}

func (r *Repository) Close() error {
	if r.conn != nil {
		slog.Info("Закрытие соединения с базой данных")
		if err := r.conn.Close(); err != nil {
			slog.Error("Ошибка при закрытии соединения с БД", "error", err)
			return err
		}
	}
	return nil
}
