package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"sync"
)

const helpMessage =
`
GOOG - тикер акции
/uptime сколько я работаю непрерывно
/stats статистика запросов
`

type helpCommand struct {
	bot *tgbotapi.BotAPI
}

var (
	helpCommandInstance *helpCommand
	helpOnce            sync.Once
)

func NewHelpCommand(bot *tgbotapi.BotAPI) Command {
	helpOnce.Do(func() {
		helpCommandInstance = &helpCommand{bot: bot}
	})
	return helpCommandInstance
}

func (c helpCommand) Invoke(update tgbotapi.Update) {
	message := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
	if _, err := c.bot.Send(message); err != nil {
		log.Errorf("help command error: %s", err)
	}
}
