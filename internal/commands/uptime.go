package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/go-memdb"
	"github.com/kettari/shitdetector/internal/commands/storage"
	log "github.com/sirupsen/logrus"
	"sync"
)

type (
	uptimeCommand struct {
		bot           *tgbotapi.BotAPI
		db            *memdb.MemDB
		uptimeService UptimeService
	}
	UptimeService interface {
		Update() error
		Since() (string, error)
	}
)

var (
	uptimeCommandInstance *uptimeCommand
	uptimeOnce            sync.Once
)

func NewUptimeCommand(bot *tgbotapi.BotAPI, db *memdb.MemDB) Command {
	uptimeOnce.Do(func() {
		us := storage.NewUptimeService(db)
		if err := us.Update(); err != nil {
			log.Error(err)
		}
		uptimeCommandInstance = &uptimeCommand{bot: bot, db: db, uptimeService: us}
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
