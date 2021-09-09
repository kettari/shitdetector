package asset

import (
	"time"
)

type (
	Stock struct {
		ID        string
		Ticker    string
		ShortName string
		Created   int64
		MarketCap float64
		EPS       float64
		ROE       float64
		Leverage  float64
		EPSRate   float64
		Currency  string
	}
	Service interface {
		Create(stock *Stock) error
		Get(ticker string) (*Stock, error)
		Delete(stock *Stock) error
	}
)

func (s Stock) Expired() bool {
	then := time.Unix(s.Created, 0)
	return time.Since(then) > time.Hour * 24
}

func (s *Stock) ExchangeToUSD(rate float64) {
	if rate == float64(0) {
		return
	}
	s.MarketCap = s.MarketCap / rate
	s.EPS = s.EPS / rate
	s.Currency = "USD"
}
