package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kettari/shitdetector/errors"
	"github.com/kettari/shitdetector/internal/asset"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

const (
	billion             = 1000000000
	stockMessagePattern = `<b>%s %s</b>
1) Рыночная капитализация (market cap), млрд: %.2f
2) Прибыльность (EPS ttm): %.2f
3) Рентабельность капитала (Return on Equity): %.2f%%
4) Леверидж (Debt/Equity): %.2f
5) Темпы роста EPS: --

<i>Актуально на %s</i>`
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
	ticker := strings.ToUpper(strings.Trim(update.Message.CommandArguments(), " "))
	if len(ticker) == 0 {
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Пустой тикер")
		if _, err := c.bot.Send(message); err != nil {
			log.Errorf("stock command error: %s", err)
			return
		}
		return
	}
	messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Загружаю акцию %s…", ticker))
	message, err := c.bot.Send(messageConfig)
	if err != nil {
		log.Errorf("stock command error: %s", err)
		return
	}

	stock, err := c.assetService.Get(ticker)
	if err != nil {
		if err == errors.ErrStockNotFound || err == errors.ErrQuoteSummaryEmpty {
			log.Infof("stock %s not found ", ticker)
			editMessageConfig := tgbotapi.NewEditMessageText(update.Message.Chat.ID, message.MessageID, "Акция не найдена")
			if _, err := c.bot.Send(editMessageConfig); err != nil {
				log.Errorf("stock command error: %s", err)
				return
			}
			return
		}

		log.Errorf("can't get stock: %s", err)

		editMessageConfig := tgbotapi.NewEditMessageText(update.Message.Chat.ID, message.MessageID, "Что-то сломалось, извините")
		if _, err := c.bot.Send(editMessageConfig); err != nil {
			log.Errorf("stock command error: %s", err)
			return
		}
		return
	}

	if stock != nil {
		editMessageConfig := tgbotapi.NewEditMessageText(update.Message.Chat.ID, message.MessageID,
			fmt.Sprintf(
				stockMessagePattern,
				ticker,
				stock.ShortName,
				stock.MarketCap/billion,
				stock.EPS,
				stock.ROE*100,
				stock.Leverage/100,
				time.Unix(stock.Created, 0).Format("Jan 2 15:04:05 2006 MST")))
		editMessageConfig.ParseMode = "HTML"
		if _, err := c.bot.Send(editMessageConfig); err != nil {
			log.Errorf("stock command error: %s", err)
			return
		}
	}
}
