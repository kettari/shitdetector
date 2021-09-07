package storage

import (
	"fmt"
	"github.com/hako/durafmt"
	"github.com/hashicorp/go-memdb"
	"time"
)

type (
	Uptime struct {
		ID    int64
		Since int64
	}
	uptimeService struct {
		db *memdb.MemDB
	}
)

func NewUptimeSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"uptime": {
				Name: "uptime",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "ID"},
					},
				},
			},
		},
	}
}

func NewUptimeService(db *memdb.MemDB) *uptimeService {
	return &uptimeService{db: db}
}

func (s uptimeService) Update() (err error) {
	txn := s.db.Txn(true)
	defer txn.Abort()

	uptime := &Uptime{ID: 1, Since: time.Now().Unix()}
	if err = txn.Insert("uptime", uptime); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s uptimeService) Since() (string, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("uptime", "id")
	if err != nil {
		return "", err
	}
	uptime, ok := raw.(*Uptime)
	if !ok {
		return "", fmt.Errorf("can't cast Uptime: %s", err)
	}

	return uptime.ToWording(), nil
}

func (u Uptime) ToWording() (wording string) {
	t := time.Unix(u.Since, 0)
	return durafmt.Parse(time.Since(t).Round(time.Second)).String()
}
