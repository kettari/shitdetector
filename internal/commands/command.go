package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Command interface {
	Invoke(tgbotapi.Update)
}
