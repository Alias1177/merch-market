package repositories

import (
	"context"
	"github.com/Alias1177/merch-store/internal/models"
)

func (u *Repository) GetUserInfo(ctx context.Context, userID int) (*models.InfoResponse, error) {
	tx, err := u.conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var coins int
	err = tx.GetContext(ctx, &coins, "SELECT coins FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, err
	}

	var inventory []models.InventoryItem
	err = tx.SelectContext(ctx, &inventory, `
        SELECT i.name, inv.quantity 
        FROM inventory inv
        JOIN items i ON inv.item_id = i.id 
        WHERE inv.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	var received []models.ReceivedTransaction
	err = tx.SelectContext(ctx, &received, `
        SELECT u.username, t.amount 
        FROM transactions t 
        JOIN users u ON t.sender_id = u.id 
        WHERE t.receiver_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	var sent []models.SentTransaction
	err = tx.SelectContext(ctx, &sent, `
        SELECT u.username, t.amount 
        FROM transactions t 
        JOIN users u ON t.receiver_id = u.id 
        WHERE t.sender_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &models.InfoResponse{
		Coins:     coins,
		Inventory: inventory,
		CoinHistory: models.CoinHistoryDetails{
			Received: received,
			Sent:     sent,
		},
	}, nil
}
