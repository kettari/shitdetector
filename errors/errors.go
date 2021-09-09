package errors

import "errors"

var (
	ErrContainerBot              = errors.New("can't get BotAPI from the container")
	ErrContainerDb               = errors.New("can't get DB from the container")
	ErrBotTokenIsEmpty           = errors.New("bot: token is empty")
	ErrStockNotFound             = errors.New("db: stock not found")
	ErrQuoteSummaryEmpty         = errors.New("yahoo provider: quote summary is empty")
	ErrCurrencyNotFound          = errors.New("currency provider: currency not found")
	ErrExchangerateApiKeyIsEmpty = errors.New("exchangerate: api key is empty")
)
