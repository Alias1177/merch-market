package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

func (r *Repository) SendCoins(ctx context.Context, senderID int, receiverUsername string, amount int) error {
	tx, err := r.conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				slog.Error("Transaction rollback failed",
					"error", rbErr,
					"original_error", err)
			}
		}
	}()

	// Проверяем существование получателя
	var receiverID int
	err = tx.GetContext(ctx, &receiverID,
		"SELECT id FROM users WHERE username = $1 FOR UPDATE",
		receiverUsername)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("user not found")
	}
	if err != nil {
		return fmt.Errorf("failed to get receiver: %w", err)
	}

	// Проверяем баланс отправителя
	var senderCoins int
	err = tx.GetContext(ctx, &senderCoins,
		"SELECT coins FROM users WHERE id = $1 FOR UPDATE",
		senderID)
	if err != nil {
		return fmt.Errorf("failed to get sender balance: %w", err)
	}

	if senderCoins < amount {
		return fmt.Errorf("not enough coins")
	}

	// Обновляем балансы
	_, err = tx.ExecContext(ctx,
		"UPDATE users SET coins = coins - $1 WHERE id = $2",
		amount, senderID)
	if err != nil {
		return fmt.Errorf("failed to update sender balance: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE users SET coins = coins + $1 WHERE id = $2",
		amount, receiverID)
	if err != nil {
		return fmt.Errorf("failed to update receiver balance: %w", err)
	}

	// Записываем транзакцию
	_, err = tx.ExecContext(ctx,
		`INSERT INTO transactions (sender_id, receiver_id, amount)
         VALUES ($1, $2, $3)`,
		senderID, receiverID, amount)
	if err != nil {
		return fmt.Errorf("failed to record transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
