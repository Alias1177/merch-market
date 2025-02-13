package models

type InfoResponse struct {
	Coins       int                `json:"coins"`
	Inventory   []InventoryItem    `json:"inventory"`
	CoinHistory CoinHistoryDetails `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type" db:"name"`
	Quantity int    `json:"quantity" db:"quantity"`
}

type CoinHistoryDetails struct {
	Received []ReceivedTransaction `json:"received"`
	Sent     []SentTransaction     `json:"sent"`
}

type ReceivedTransaction struct {
	FromUser string `json:"fromUser" db:"username"`
	Amount   int    `json:"amount" db:"amount"`
}

type SentTransaction struct {
	ToUser string `json:"toUser" db:"username"`
	Amount int    `json:"amount" db:"amount"`
}
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
