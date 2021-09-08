package storage

import (
	"fmt"
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/go-uuid"
	"github.com/kettari/shitdetector/internal/stock_log"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

const maxEntriesCount = 10000

type (
	stockLogService struct {
		db *memdb.MemDB
	}
)

func NewStockLogSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"stock_log": {
				Name: "stock_log",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UUIDFieldIndex{Field: "ID"},
					},
					"timestamp": {
						Name:    "timestamp",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Created"},
					},
					"Ticker": {
						Name:    "Ticker",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Ticker"},
					},
				},
			},
		},
	}
}

func NewStockLogService(db *memdb.MemDB) *stockLogService {
	return &stockLogService{db: db}
}

func (s stockLogService) Log(ticker string) (err error) {
	txn := s.db.Txn(true)
	defer txn.Abort()

	id, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}

	stlog := &stock_log.StockLog{
		ID:      id,
		Created: time.Now().Unix(),
		Ticker:  ticker,
	}
	if err = txn.Insert("stock_log", stlog); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s stockLogService) Cleanup() (err error) {
	txn := s.db.Txn(true)
	defer txn.Abort()

	it, err := txn.GetReverse("stock_log", "timestamp")
	if err != nil {
		return err
	}
	count := 0
	cleaned := 0
	for obj := it.Next(); obj != nil; obj = it.Next() {
		stockLog, ok := obj.(*stock_log.StockLog)
		if !ok {
			return fmt.Errorf("can't cast StockLog: %s", err)
		}
		count++
		if count > maxEntriesCount {
			cleaned++
			if err = txn.Delete("stock_log", stockLog); err != nil {
				return err
			}
		}
		log.Debugf("stock_log item %d: timestamp=%d, id=%s, Ticker=%s", count, stockLog.Created, stockLog.ID, stockLog.Ticker)
	}
	log.Infof("stock_log cleanup, total %d records found, from which %d cleaned", count, cleaned)

	txn.Commit()
	return nil
}

func (s stockLogService) Stats() (stats stock_log.TickerStats, err error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("stock_log", "timestamp")
	if err != nil {
		return nil, err
	}
	count := 0
	var tickers = map[string]int64{}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		stockLog, ok := obj.(*stock_log.StockLog)
		if !ok {
			return nil, fmt.Errorf("can't cast StockLog: %s", err)
		}
		count++
		if _, ok := tickers[stockLog.Ticker]; !ok {
			tickers[stockLog.Ticker] = 1
		} else {
			tickers[stockLog.Ticker]++
		}
	}
	log.Debugf("stock_log stats based on %d records", count)

	stats = make(stock_log.TickerStats, 0)
	for tckr, cnt := range tickers {
		t := &stock_log.TickerStat{
			Ticker: tckr,
			Count:  cnt,
		}
		stats = append(stats, t)
	}
	sort.Sort(stats)

	return stats, nil
}
