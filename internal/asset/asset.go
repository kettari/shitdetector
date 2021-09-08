package asset

import "time"

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
	}
	Service interface {
		Create(stock *Stock) error
		Get(ticker string) (*Stock, error)
		Update(stock *Stock) error
		Delete(stock *Stock) error
	}
)

func (s Stock) Expired() bool {
	then := time.Unix(s.Created, 0)
	return time.Since(then) > time.Hour
}
