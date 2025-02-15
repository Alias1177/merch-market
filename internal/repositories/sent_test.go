package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendCoins(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		mock.ExpectBegin()

		// Проверка существования получателя
		mock.ExpectQuery("SELECT id FROM users WHERE username = \\$1").
			WithArgs("receiver").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		// Проверка баланса отправителя
		mock.ExpectQuery("SELECT coins FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(1000))

		// Обновление баланса отправителя
		mock.ExpectExec("UPDATE users SET coins = coins - \\$1 WHERE id = \\$2").
			WithArgs(500, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Обновление баланса получателя
		mock.ExpectExec("UPDATE users SET coins = coins \\+ \\$1 WHERE id = \\$2").
			WithArgs(500, 2).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Запись транзакции
		mock.ExpectExec("INSERT INTO transactions \\(sender_id, receiver_id, amount\\)").
			WithArgs(1, 2, 500).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = repo.SendCoins(context.Background(), 1, "receiver", 500)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("receiver not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT id FROM users WHERE username = \\$1").
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		mock.ExpectRollback()

		err = repo.SendCoins(context.Background(), 1, "nonexistent", 500)
		assert.EqualError(t, err, "user not found")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("not enough coins", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT id FROM users WHERE username = \\$1").
			WithArgs("receiver").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		mock.ExpectQuery("SELECT coins FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(50))

		// Перемещаем проверку ошибки после ExpectationsWereMet
		err = repo.SendCoins(context.Background(), 1, "receiver", 100)

		// Даем возможность выполниться rollback до проверки ошибки
		errExpectations := mock.ExpectationsWereMet()
		assert.NoError(t, errExpectations)

		// Теперь проверяем ошибку отправки монет
		assert.EqualError(t, err, "not enough coins")
	})

	t.Run("transaction begin error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		err = repo.SendCoins(context.Background(), 1, "receiver", 500)
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("update sender balance error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT id FROM users WHERE username = \\$1").
			WithArgs("receiver").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		mock.ExpectQuery("SELECT coins FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(1000))

		mock.ExpectExec("UPDATE users SET coins = coins - \\$1 WHERE id = \\$2").
			WithArgs(500, 1).
			WillReturnError(sql.ErrConnDone)

		mock.ExpectRollback()

		err = repo.SendCoins(context.Background(), 1, "receiver", 500)
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
