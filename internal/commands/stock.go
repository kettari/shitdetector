package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kettari/shitdetector/errors"
	"github.com/kettari/shitdetector/internal/asset"
	"github.com/kettari/shitdetector/internal/stock_log"
	"github.com/kettari/shitdetector/internal/underwriter/finindie"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	billion             = 1000000000
	stockMessagePattern = `<b>%s %s</b> %s
1) Рыночная капитализация (market cap): $%.2f млрд
2) Прибыльность (EPS ttm): $%.2f
3) Рентабельность капитала (ROE): %.2f%%
4) Леверидж (Debt/Equity): %.2f
5) Темпы роста EPS: %.2f%%

<i>Актуально на %s</i>`
	badScore    = "\xF0\x9F\x92\xA9"
	normalScore = "\xF0\x9F\x92\xAA"
	goodScore   = "\xF0\x9F\x9A\x80"
)

type (
	stockCommand struct {
		bot             *tgbotapi.BotAPI
		assetService    asset.Service
		stockLogService stock_log.Service
	}
)

var (
	stockCommandInstance *stockCommand
	stockOnce            sync.Once
	suggestions          = map[string]string{
		"APPL": "AAPL -- Apple Inc.",
		"LKOH": "LKOH.ME -- OIL CO LUKOIL PJSC",
		"SBER": "SBER.ME -- SBERBANK OF RUSSIA",
		"ENRU": "ENEL.ME -- ENEL RUSSIA PJSC",
		"RUAL": "RUAL.ME -- UNITED COMPANY RU",
		"TCSG": "TCSG.ME -- TCS GROUP HOLDING",
		"GAZP": "GAZP.ME -- GAZPROM PJSC",
		"FIXP": "FIXP.ME -- FIX PRICE GROUP LT",
		"IRAO": "IRAO.ME -- INTER RAO UES PJSC",
		"CHMF": "CHMF.ME -- SEVERSTAL PJSC",
		"POLY": "POLY.ME -- POLYMETAL INTL PLC",
		"POGR": "POGR.ME -- PETROPAVLOVSK PLC",
		"ALRS": "ALRS.ME -- ALROSA PJSC",
		"AFLT": "AFLT.ME -- AEROFLOT PJSC",
	}
)

func NewStockCommand(bot *tgbotapi.BotAPI, assetSvc asset.Service, stockLogSvc stock_log.Service) Command {
	stockOnce.Do(func() {
		stockCommandInstance = &stockCommand{bot: bot, assetService: assetSvc, stockLogService: stockLogSvc}
	})
	return stockCommandInstance
}

func (c stockCommand) Invoke(update tgbotapi.Update) {
	ticker := ""
	if update.Message.IsCommand() {
		ticker = update.Message.CommandArguments()
	} else {
		ticker = update.Message.Text
	}
	ticker = strings.ToUpper(strings.Trim(ticker, " "))
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

	if sugg, ok := suggestions[ticker]; ok {
		messageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Возможно, вы имели в виду %s", sugg))
		_, err := c.bot.Send(messageConfig)
		if err != nil {
			log.Errorf("stock command error: %s", err)
			return
		}
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
		underwriter := finindie.NewFinindieUnderwriter(stock)
		score := underwriter.Score()
		totalScore := score.TotalScore()

		scoreSign := normalScore
		if totalScore < 15 {
			scoreSign = badScore
		} else if totalScore > 20 {
			scoreSign = goodScore
		}

		text := fmt.Sprintf(
			stockMessagePattern,
			ticker,
			stock.ShortName,
			scoreSign,
			stock.MarketCap/billion,
			stock.EPS,
			stock.ROE*100,
			stock.Leverage/100,
			stock.EPSRate*100,
			time.Unix(stock.Created, 0).Format("Jan 2 15:04:05 2006 MST"))

		text += "\n\n" + score.Describe() + "\n\n<i>Не является индивидуальной инвестиционной рекомендацией</i>"

		editMessageConfig := tgbotapi.NewEditMessageText(update.Message.Chat.ID, message.MessageID, text)
		editMessageConfig.ParseMode = "HTML"
		editMessageConfig.DisableWebPagePreview = true
		if _, err := c.bot.Send(editMessageConfig); err != nil {
			log.Errorf("stock command error: %s", err)
			return
		}

		// Log request
		if err = c.stockLogService.Log(stock.Ticker); err != nil {
			log.Errorf("stock command error: %s", err)
			return
		}
		// Maintenance
		probability := randInt(1, 100)
		if probability > 80 {
			log.Info("performing stock_log cleanup")
			if err = c.stockLogService.Cleanup(); err != nil {
				log.Errorf("stock command error: %s", err)
				return
			}
		} else {
			log.Debug("skipping stock_log cleanup")
		}
	}
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
