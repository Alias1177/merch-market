package pkg

import "errors"

var (
	DbError              = "Error connecting to the database ⬇️"
	CfgErr               = "Error reading config file:⬇️"
	ParseCfgErr          = "Error parsing config:⬇️"
	ErrUserAlreadyExists = errors.New("user already exists")
)
