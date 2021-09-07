package errors

import "errors"

var (
	ErrBotTokenIsEmpty = errors.New("bot: token is empty")
	ErrBotSendMessage = errors.New("bot: send message error: ")
)
