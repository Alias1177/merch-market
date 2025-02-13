package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func (r *Repository) SendCoins(ctx context.Context, senderID int, receiverUsername string, amount int) error {
	tx, err := r.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Ошибка при откате транзакции: %v", rbErr)
			}
		}
	}()

	var receiverID int
	err = tx.GetContext(ctx, &receiverID, "SELECT id FROM users WHERE username = $1", receiverUsername)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("user not found")
	} else if err != nil {
		return err
	}

	var senderCoins int
	err = tx.GetContext(ctx, &senderCoins, "SELECT coins FROM users WHERE id = $1", senderID)
	if err != nil {
		return err
	}

	if senderCoins < amount {
		return fmt.Errorf("not enough coins")
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET coins = coins - $1 WHERE id = $2", amount, senderID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET coins = coins + $1 WHERE id = $2", amount, receiverID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
        INSERT INTO transactions (sender_id, receiver_id, amount)
        VALUES ($1, $2, $3)
    `, senderID, receiverID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}
