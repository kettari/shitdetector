package provider

import "github.com/kettari/shitdetector/internal/asset"

type StockProvider interface {
	Fetch(ticker string) (*asset.Stock, error)
}
