package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kettari/shitdetector/internal/stock_log"
	log "github.com/sirupsen/logrus"
	"sync"
)

type (
	statsCommand struct {
		bot             *tgbotapi.BotAPI
		stockLogService stock_log.Service
	}
)

var (
	statsCommandInstance *statsCommand
	stockLogOnce         sync.Once
)

func NewStatsCommand(bot *tgbotapi.BotAPI, service stock_log.Service) Command {
	stockLogOnce.Do(func() {
		statsCommandInstance = &statsCommand{bot: bot, stockLogService: service}
	})
	return statsCommandInstance
}

func (c statsCommand) Invoke(update tgbotapi.Update) {
	stats, err := c.stockLogService.Stats()
	if err != nil {
		log.Errorf("can't get Stats(): %s", err)
		return
	}
	text := "<b>Статистика запросов:</b>\n"
	place := 1
	for _, v := range stats {
		text += fmt.Sprintf("%d) %s - %d\n", place, v.Ticker, v.Count)
		place++
		if place == 11 {
			break
		}
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	message.ParseMode = "HTML"
	if _, err := c.bot.Send(message); err != nil {
		log.Errorf("stats command error: %s", err)
	}
}
