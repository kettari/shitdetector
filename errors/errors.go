package errors

import "errors"

var (
	ErrContainerBot    = errors.New("can't get BotAPI from the container")
	ErrContainerDb     = errors.New("can't get DB from the container")
	ErrBotTokenIsEmpty = errors.New("bot: token is empty")
)
