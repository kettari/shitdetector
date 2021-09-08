package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"sync"
)

const sourceMessage =
`
Информация загружается из сервиса <a href="https://finance.yahoo.com">Yahoo Finance</a>, кешируется на 1 час.
`

type sourceCommand struct {
	bot *tgbotapi.BotAPI
}

var (
	sourceCommandInstance *sourceCommand
	sourceOnce            sync.Once
)

func NewSourceCommand(bot *tgbotapi.BotAPI) Command {
	sourceOnce.Do(func() {
		sourceCommandInstance = &sourceCommand{bot: bot}
	})
	return sourceCommandInstance
}

func (c sourceCommand) Invoke(update tgbotapi.Update) {
	message := tgbotapi.NewMessage(update.Message.Chat.ID, sourceMessage)
	message.ParseMode = "HTML"
	if _, err := c.bot.Send(message); err != nil {
		log.Errorf("help command error: %s", err)
	}
}
