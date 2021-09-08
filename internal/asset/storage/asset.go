package storage

import (
	"fmt"
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/go-uuid"
	"github.com/kettari/shitdetector/errors"
	"github.com/kettari/shitdetector/internal/asset"
	"github.com/kettari/shitdetector/internal/provider"
	log "github.com/sirupsen/logrus"
)

type store struct {
	db   *memdb.MemDB
	prov provider.StockProvider
}

func NewAssetService(db *memdb.MemDB, prov provider.StockProvider) asset.Service {
	return &store{db: db, prov: prov}
}

func NewAssetSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"asset": {
				Name: "asset",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UUIDFieldIndex{Field: "ID"},
					},
					"ticker": {
						Name:    "ticker",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Ticker"},
					},
				},
			},
		},
	}
}

func (s store) Create(stock *asset.Stock) error {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	stock.ID = id
	return s.Update(stock)
}

func (s store) Update(stock *asset.Stock) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert("asset", stock); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s store) Get(ticker string) (*asset.Stock, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	var stock *asset.Stock
	shouldFetch := false
	raw, err := txn.First("asset", "ticker", ticker)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		shouldFetch = true
		log.Debugf("stock %s is not found in the db, should fetch", ticker)
	} else {
		var ok bool
		stock, ok = raw.(*asset.Stock)
		if !ok {
			return nil, fmt.Errorf("can't cast Asset: %s", err)
		}
		log.Debugf("stock %s is found in the db", ticker)
		if stock.Expired() {
			shouldFetch = true
			log.Debugf("stock %s expired, should fetch", ticker)
		}
	}

	if shouldFetch {
		fetchedStock, err := s.prov.Fetch(ticker)
		if err != nil {
			return nil, err
		}
		if fetchedStock == nil {
			return nil, errors.ErrStockNotFound
		}
		log.Debugf("stock %s fetched from provider", ticker)
		if stock != nil {
			if err = s.Update(fetchedStock); err != nil {
				return nil, err
			}
		} else {
			if err = s.Create(fetchedStock); err != nil {
				return nil, err
			}
			stock = fetchedStock
		}
	}

	return stock, nil
}
