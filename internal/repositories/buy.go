package repositories

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// BuyRepository реализует интерфейс BuyRepo
type BuyRepository struct {
	conn *sqlx.DB
}

func NewBuyRepository(repo *Repository) *BuyRepository {
	return &BuyRepository{conn: repo.conn}
}

// Реализация метода BuyItem (выполнение транзакции)
func (r *BuyRepository) BuyItem(ctx context.Context, userID, itemID int) error {
	tx, err := r.conn.BeginTxx(ctx, nil) // Начинаем транзакцию
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // Если случилась паника, откатываем транзакцию
		} else if err != nil {
			_ = tx.Rollback() // Откатываем транзакцию, если была ошибка
		} else {
			_ = tx.Commit() // Если ошибок нет, фиксируем транзакцию
		}
	}()

	var price int
	if err = tx.GetContext(ctx, &price, "SELECT price FROM items WHERE id = $1", itemID); err != nil {
		return fmt.Errorf("item not found or failed to get price: %w", err)
	}

	var coins int
	if err = tx.GetContext(ctx, &coins, "SELECT coins FROM users WHERE id = $1", userID); err != nil {
		return fmt.Errorf("user not found or failed to get balance: %w", err)
	}

	if coins < price {
		return fmt.Errorf("not enough coins for the purchase")
	}

	if _, err = tx.ExecContext(ctx, "UPDATE users SET coins = coins - $1 WHERE id = $2", price, userID); err != nil {
		return fmt.Errorf("failed to update user coins: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO inventory (user_id, item_id, quantity)
		VALUES ($1, $2, 1)
		ON CONFLICT (user_id, item_id)
		DO UPDATE SET quantity = inventory.quantity + 1
	`, userID, itemID)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	return nil
}
