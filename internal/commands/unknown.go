package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"sync"
)

type unknownCommand struct {
	bot *tgbotapi.BotAPI
}

var (
	unknownCommandInstance *unknownCommand
	unknownOnce            sync.Once
)

func NewUnknownCommand(bot *tgbotapi.BotAPI) Command {
	unknownOnce.Do(func() {
		unknownCommandInstance = &unknownCommand{bot: bot}
	})
	return unknownCommandInstance
}

func (c unknownCommand) Invoke(update tgbotapi.Update) {
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Незнакомая команда. Попробуйте /help")
	if _, err := c.bot.Send(message); err != nil {
		log.Errorf("unknown command error: %s", err)
	}
}
