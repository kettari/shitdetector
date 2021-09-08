package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kettari/shitdetector/errors"
	"github.com/kettari/shitdetector/internal/asset"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type (
	stockCommand struct {
		bot          *tgbotapi.BotAPI
		assetService asset.Service
	}
)

var (
	stockCommandInstance *stockCommand
	stockOnce            sync.Once
)

func NewStockCommand(bot *tgbotapi.BotAPI, service asset.Service) Command {
	stockOnce.Do(func() {
		stockCommandInstance = &stockCommand{bot: bot, assetService: service}
	})
	return stockCommandInstance
}

func (c stockCommand) Invoke(update tgbotapi.Update) {
	ticker := strings.Trim(update.Message.CommandArguments(), " ")
	if len(ticker) == 0 {
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Пустой тикер")
		if _, err := c.bot.Send(message); err != nil {
			log.Errorf("stock command error: %s", err)
			return
		}
		return
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Извлекаю акцию %s", ticker))
	if _, err := c.bot.Send(message); err != nil {
		log.Errorf("stock command error: %s", err)
		return
	}

	stock, err := c.assetService.Get(ticker)
	if err != nil {
		if err == errors.ErrStockNotFound {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Акция не найдена")
			if _, err := c.bot.Send(message); err != nil {
				log.Errorf("stock command error: %s", err)
				return
			}
		}
		log.Errorf("can't get stock: %s", err)
		return
	}

	if stock != nil {
		message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("ID stock %s", stock.ID))
		if _, err := c.bot.Send(message); err != nil {
			log.Errorf("stock command error: %s", err)
			return
		}
	}
}
