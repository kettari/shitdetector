package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kettari/shitdetector/internal/uptime"
	log "github.com/sirupsen/logrus"
	"sync"
)

type (
	uptimeCommand struct {
		bot           *tgbotapi.BotAPI
		uptimeService uptime.Service
	}
)

var (
	uptimeCommandInstance *uptimeCommand
	uptimeOnce            sync.Once
)

func NewUptimeCommand(bot *tgbotapi.BotAPI, service uptime.Service) Command {
	uptimeOnce.Do(func() {
		uptimeCommandInstance = &uptimeCommand{bot: bot, uptimeService: service}
	})
	return uptimeCommandInstance
}

func (c uptimeCommand) Invoke(update tgbotapi.Update) {
	since, err := c.uptimeService.Since()
	if err != nil {
		log.Errorf("can't get uptime Since(): %s", err)
		return
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Работаю %s", since))
	if _, err := c.bot.Send(message); err != nil {
		log.Errorf("uptime command error: %s", err)
	}
}
