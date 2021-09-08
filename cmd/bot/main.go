package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/go-memdb"
	"github.com/kettari/shitdetector/errors"
	storage2 "github.com/kettari/shitdetector/internal/asset/storage"
	"github.com/kettari/shitdetector/internal/commands"
	"github.com/kettari/shitdetector/internal/config"
	"github.com/kettari/shitdetector/internal/provider/yahoo"
	"github.com/kettari/shitdetector/internal/registry"
	"github.com/kettari/shitdetector/internal/uptime/storage"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	cnt, err := registry.NewContainer(config.GetConfig())
	if err != nil {
		log.Panic(err)
	}
	bot, ok := cnt.Get("bot").(*tgbotapi.BotAPI)
	if !ok {
		log.Panic(errors.ErrContainerBot)
	}
	log.Infof("authorized on account %s", bot.Self.UserName)

	db, ok := cnt.Get("db").(*memdb.MemDB)
	if !ok {
		log.Panic(errors.ErrContainerDb)
	}
	uptimeService := storage.NewUptimeService(db)
	if err := uptimeService.Update(); err != nil {
		log.Panic(err)
	}
	assetService := storage2.NewAssetService(db, yahoo.NewYahooProvider())

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	// Wait for updates and clear them as we don't want to handle a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if len(update.Message.Command()) > 0 {
			log.Infof("[%s] command: %s", update.Message.From.UserName, update.Message.Text)

			var cmd commands.Command
			switch update.Message.Command() {
			case "help":
				cmd = commands.NewHelpCommand(bot)
			case "stock":
				cmd = commands.NewStockCommand(bot, assetService)
			case "uptime":
				cmd = commands.NewUptimeCommand(bot, uptimeService)
			default:
				cmd = commands.NewUnknownCommand(bot)
			}

			go cmd.Invoke(update)

		} else {
			log.Infof("[%s] message: %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Непонятно. Пришлите мне тикер или команду, пожалуйста /help")
			msg.ReplyToMessageID = update.Message.MessageID

			_, err = bot.Send(msg)
			if err != nil {
				log.Error(err)
				break
			}
		}
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)

	err := tgbotapi.SetLogger(log.StandardLogger())
	if err != nil {
		log.Panic(err)
	}
}
