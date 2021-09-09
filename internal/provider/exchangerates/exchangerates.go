package exchangerates

import (
	"encoding/json"
	"fmt"
	"github.com/kettari/shitdetector/internal/currency"
	"net/http"
	"time"
)

const exchangeratesUrl = `http://api.exchangeratesapi.io/v1/latest?access_key=%s&base=EUR&symbols=USD,%s`

type (
	exchangeratesProvider struct {
		apiKey string
	}
	exchangeRatesResponse struct {
		Success   bool               `json:"success"`
		Timestamp int64              `json:"timestamp"`
		Base      string             `json:"base"`
		Date      string             `json:"date"`
		Rates     map[string]float64 `json:"rates"`
	}
)

func NewExchangeratesProvider(apiKey string) *exchangeratesProvider {
	return &exchangeratesProvider{
		apiKey: apiKey,
	}
}

func (p exchangeratesProvider) Fetch(symbol string) (curr *currency.Currency, err error) {
	resp, err := http.Get(fmt.Sprintf(exchangeratesUrl, p.apiKey, symbol))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	exchResp := &exchangeRatesResponse{}
	if err = json.NewDecoder(resp.Body).Decode(exchResp); err != nil {
		return nil, err
	}

	symbolToUSDRate := exchResp.Rates[symbol] / exchResp.Rates["USD"]

	return &currency.Currency{
		Symbol:  symbol,
		Rate:    symbolToUSDRate,
		Created: time.Now().Unix(),
	}, nil
}
