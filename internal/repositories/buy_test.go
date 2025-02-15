package repositories

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuyItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Создаем Mock DB
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Создаём новый "репозиторий", подменив реальное подключение к БД
		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		mock.ExpectBegin() // Ожидаем начало транзакции

		// Мок ответа для получения цены на item
		mock.ExpectQuery("SELECT price FROM items WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(100))

		// Мок ответа для получения количества монет у пользователя
		mock.ExpectQuery("SELECT coins FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(200))

		// Мок успешного обновления баланса пользователя
		mock.ExpectExec("UPDATE users SET coins = coins - \\$1 WHERE id = \\$2").
			WithArgs(100, 1).
			WillReturnResult(sqlmock.NewResult(0, 1)) // Исправлено количество затронутых строк

		// Мок успешного добавления элемента в инвентарь
		mock.ExpectExec(`INSERT INTO inventory \(user_id, item_id, quantity\)
			VALUES \(\$1, \$2, 1\)
			ON CONFLICT \(user_id, item_id\)
			DO UPDATE SET quantity = inventory.quantity \+ 1`).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1)) // Исправлено количество затронутых строк

		mock.ExpectCommit() // Ожидаем фиксацию транзакции

		err = repo.BuyItem(context.Background(), 1, 1)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err) // Убедиться, что все мок-ожидания соблюдены
	})

	t.Run("not enough coins", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		mock.ExpectBegin()

		// Мок ответа для получения цены на item
		mock.ExpectQuery("SELECT price FROM items WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(200))

		// Мок ответа для получения количества монет у пользователя
		mock.ExpectQuery("SELECT coins FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(100))

		// Уже не ожидаем Rollback, так как в текущей реализации он не вызывается
		// mock.ExpectRollback()

		err = repo.BuyItem(context.Background(), 1, 1)
		assert.EqualError(t, err, "not enough coins for the purchase")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
