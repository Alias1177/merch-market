package pkg

import "errors"

var (
	DbError              = "Error connecting to the database ⬇️"
	CfgErr               = "Error reading config file:⬇️"
	ParseCfgErr          = "Error parsing config:⬇️"
	ErrItemNotFound      = errors.New("Item not found")
	ErrUserNotFound      = errors.New("User not found")
	ErrInsufficientFunds = errors.New("Insufficient funds")
	ErrUserAlreadyExists = errors.New("user already exists")
)
