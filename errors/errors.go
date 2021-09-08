package errors

import "errors"

var (
	ErrContainerBot       = errors.New("can't get BotAPI from the container")
	ErrContainerDb        = errors.New("can't get DB from the container")
	ErrBotTokenIsEmpty    = errors.New("bot: token is empty")
	ErrStockAlreadyExists = errors.New("db: stock with this ticker already exists")
	ErrStockNotFound      = errors.New("db: stock not found")
)
