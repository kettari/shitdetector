package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"sync"
)

const helpMessage =
`
Рассчитываю скоринг по методу <a href="https://t.me/Finindie/767">Finindie</a>, показываю основные параметры акций.

Просто прислать текст «GOOG» - тикер акции
/uptime сколько я работаю непрерывно
/stats статистика запросов
/stock информация по акции
/source информация по источникам информации

<i>Информация от бота не является персональной инвестиционной рекомендацией</i>
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
	message.ParseMode = "HTML"
	message.DisableWebPagePreview = true
	if _, err := c.bot.Send(message); err != nil {
		log.Errorf("help command error: %s", err)
	}
}
