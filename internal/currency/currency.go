package currency

import "time"

type (
	Currency struct {
		ID        string
		Symbol    string
		Rate      float64
		Created   int64
	}
	Service interface {
		Create(currency *Currency) error
		Get(symbol string) (*Currency, error)
		Delete(currency *Currency) error
	}
)

func (s Currency) Expired() bool {
	then := time.Unix(s.Created, 0)
	return time.Since(then) > time.Hour * 48
}
