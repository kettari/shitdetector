package provider

import (
	"github.com/kettari/shitdetector/internal/asset"
	"github.com/kettari/shitdetector/internal/currency"
)

type (
	StockProvider interface {
		Fetch(ticker string) (*asset.Stock, error)
	}
	CurrencyProvider interface {
		Fetch(symbol string) (*currency.Currency, error)
	}
)
