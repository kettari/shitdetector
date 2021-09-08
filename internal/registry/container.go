package registry

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/go-memdb"
	"github.com/kettari/shitdetector/internal/asset/storage"
	"github.com/kettari/shitdetector/internal/config"
	storage2 "github.com/kettari/shitdetector/internal/uptime/storage"
	"github.com/sarulabs/di"
)

func NewContainer(conf *config.Config) (di.Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}

	if err := builder.Add([]di.Def{
		{
			Name: "bot",
			Build: func(ctn di.Container) (interface{}, error) {
				bot, err := tgbotapi.NewBotAPI(conf.BotToken)
				if err != nil {
					return nil, err
				}
				bot.Debug = conf.Debug
				return bot, nil
			},
		},
		{
			Name: "db",
			Build: func(cnt di.Container) (interface{}, error) {
				db, err := memdb.NewMemDB(unionSchemas(storage.NewAssetSchema(), storage2.NewUptimeSchema()))
				if err != nil {
					return nil, err
				}
				return db, nil
			},
		},
	}...); err != nil {
		return nil, err
	}

	cnt := builder.Build()

	return cnt, nil
}

func unionSchemas(schemas ...*memdb.DBSchema) *memdb.DBSchema {
	tables := map[string]*memdb.TableSchema{}
	for _, s := range schemas {
		for k, v := range s.Tables {
			tables[k] = v
		}
	}
	return &memdb.DBSchema{Tables: tables}
}
