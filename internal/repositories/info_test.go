package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Alias1177/merch-store/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		mock.ExpectBegin()

		// Мок запроса монет
		mock.ExpectQuery("SELECT coins FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(1000))

		// Мок запроса инвентаря
		mock.ExpectQuery("SELECT i.name, inv.quantity FROM inventory").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"name", "quantity"}).
				AddRow("Item1", 2).
				AddRow("Item2", 1))

		// Мок запроса полученных транзакций
		mock.ExpectQuery("SELECT u.username, t.amount FROM transactions t").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"username", "amount"}).
				AddRow("sender1", 100).
				AddRow("sender2", 200))

		// Мок запроса отправленных транзакций
		mock.ExpectQuery("SELECT u.username, t.amount FROM transactions t").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"username", "amount"}).
				AddRow("receiver1", 50).
				AddRow("receiver2", 150))

		mock.ExpectCommit()

		expectedResponse := &models.InfoResponse{
			Coins: 1000,
			Inventory: []models.InventoryItem{
				{Type: "Item1", Quantity: 2},
				{Type: "Item2", Quantity: 1},
			},
			CoinHistory: models.CoinHistoryDetails{
				Received: []models.ReceivedTransaction{
					{FromUser: "sender1", Amount: 100},
					{FromUser: "sender2", Amount: 200},
				},
				Sent: []models.SentTransaction{
					{ToUser: "receiver1", Amount: 50},
					{ToUser: "receiver2", Amount: 150},
				},
			},
		}

		info, err := repo.GetUserInfo(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, info)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("begin transaction error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		info, err := repo.GetUserInfo(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, info)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("coins query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT coins FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		info, err := repo.GetUserInfo(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, info)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("inventory query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT coins FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(1000))

		mock.ExpectQuery("SELECT i.name, inv.quantity FROM inventory").
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		mock.ExpectRollback()

		info, err := repo.GetUserInfo(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, info)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("commit error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		mock.ExpectBegin()

		// Мок запроса монет
		mock.ExpectQuery("SELECT coins FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(1000))

		// Мок запроса инвентаря
		mock.ExpectQuery("SELECT i.name, inv.quantity FROM inventory").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"name", "quantity"}))

		// Мок запроса полученных транзакций
		mock.ExpectQuery("SELECT u.username, t.amount FROM transactions t").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"username", "amount"}))

		// Мок запроса отправленных транзакций
		mock.ExpectQuery("SELECT u.username, t.amount FROM transactions t").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"username", "amount"}))

		// Ожидаем ошибку при коммите
		mock.ExpectCommit().WillReturnError(sql.ErrTxDone)

		info, err := repo.GetUserInfo(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, info)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
