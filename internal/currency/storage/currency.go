package storage

import (
	"fmt"
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/go-uuid"
	"github.com/kettari/shitdetector/errors"
	"github.com/kettari/shitdetector/internal/currency"
	"github.com/kettari/shitdetector/internal/provider"
	log "github.com/sirupsen/logrus"
)

type store struct {
	db   *memdb.MemDB
	prov provider.CurrencyProvider
}

func NewCurrencyService(db *memdb.MemDB, prov provider.CurrencyProvider) *store {
	return &store{db: db, prov: prov}
}

func NewCurrencySchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"currency": {
				Name: "currency",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UUIDFieldIndex{Field: "ID"},
					},
					"symbol": {
						Name:    "symbol",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Symbol"},
					},
				},
			},
		},
	}
}

func (s store) Create(currency *currency.Currency) error {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	currency.ID = id

	txn := s.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert("currency", currency); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s store) Delete(currency *currency.Currency) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	if err := txn.Delete("currency", currency); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s store) Get(symbol string) (*currency.Currency, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	var curr *currency.Currency
	shouldFetch := false
	raw, err := txn.First("currency", "symbol", symbol)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		shouldFetch = true
		log.Debugf("currency %s is not found in the db, should fetch", symbol)
	} else {
		var ok bool
		curr, ok = raw.(*currency.Currency)
		if !ok {
			return nil, fmt.Errorf("can't cast Currency: %s", err)
		}
		log.Debugf("currency %s is found in the db", symbol)
		if curr.Expired() {
			shouldFetch = true
			log.Debugf("currency %s expired, should fetch", symbol)
		}
	}

	if shouldFetch {
		fetchedCurr, err := s.prov.Fetch(symbol)
		if err != nil {
			return nil, err
		}
		if fetchedCurr == nil {
			return nil, errors.ErrCurrencyNotFound
		}
		log.Debugf("currency %s fetched from provider", symbol)
		if curr != nil {
			if err = s.Delete(curr); err != nil {
				return nil, err
			}
		}
		if err = s.Create(fetchedCurr); err != nil {
			return nil, err
		}
		curr = fetchedCurr
	}

	return curr, nil
}
