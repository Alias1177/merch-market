package repositories

import (
	"context"
	"database/sql"
	"github.com/Alias1177/merch-store/internal/models"
	"github.com/Alias1177/merch-store/pkg"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		username := "testuser"
		passwordHash := "hashedpassword"
		coins := 100

		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "coins"}).
			AddRow(1, username, passwordHash, coins)

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(username, passwordHash, coins).
			WillReturnRows(rows)

		expectedUser := &models.User{
			ID:           1,
			Username:     username,
			PasswordHash: passwordHash,
			Coins:        coins,
		}

		user, err := repo.CreateUser(context.Background(), username, passwordHash, coins)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("duplicate username", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		username := "existinguser"
		passwordHash := "hashedpassword"
		coins := 100

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(username, passwordHash, coins).
			WillReturnError(&pq.Error{
				Message: "duplicate key value violates unique constraint \"users_username_key\"",
			})

		user, err := repo.CreateUser(context.Background(), username, passwordHash, coins)
		assert.ErrorIs(t, err, pkg.ErrUserAlreadyExists)
		assert.Nil(t, user)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := &Repository{conn: sqlx.NewDb(db, "sqlmock")}

		username := "testuser"
		passwordHash := "hashedpassword"
		coins := 100

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(username, passwordHash, coins).
			WillReturnError(sql.ErrConnDone)

		user, err := repo.CreateUser(context.Background(), username, passwordHash, coins)
		assert.Error(t, err)
		assert.Nil(t, user)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
