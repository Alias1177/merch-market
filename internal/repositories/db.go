package repositories

import (
	"github.com/Alias1177/merch-store/pkg/logger"
	"log"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Alias1177/merch-store/pkg"
)

func Connect(dsn string) *sqlx.DB {
	logger.ColorLogger()
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		slog.Error(pkg.DbError)
		log.Fatalf("Unable to connect to the database: %v", err)
	}
	log.Println("Successfully connected to the database.")
	return db
}
