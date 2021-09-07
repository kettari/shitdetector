package storage

import (
	"github.com/hashicorp/go-memdb"
	"github.com/kettari/shitdetector/internal/asset"
)

type store struct {
	db *memdb.MemDB
}

func NewAssetService(db *memdb.MemDB) asset.Service {
	return &store{db: db}
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
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
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

func (s store) Get(ticket string) (asset *asset.Stock, err error) {
	panic("not implemented")
	return nil, nil
}
