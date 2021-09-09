package usdficator

import (
	"github.com/kettari/shitdetector/internal/asset"
	"github.com/kettari/shitdetector/internal/currency"
	"github.com/kettari/shitdetector/internal/provider"
)

type (
	usdficatorProvider struct {
		stockProvider   provider.StockProvider
		currencyService currency.Service
	}
)

func NewUSDFicatorProvider(stockProvider provider.StockProvider, currencyService currency.Service) *usdficatorProvider {
	return &usdficatorProvider{stockProvider: stockProvider, currencyService: currencyService}
}

func (p usdficatorProvider) Fetch(ticker string) (stock *asset.Stock, err error) {
	stock, err = p.stockProvider.Fetch(ticker)
	if err != nil {
		return nil, err
	}

	if stock.Currency != "USD" {
		curr, err := p.currencyService.Get(stock.Currency)
		if err != nil {
			return nil, err
		}
		stock.ExchangeToUSD(curr.Rate)
	}

	return stock, nil
}
