package cbrf

import (
	"fmt"
	"github.com/kettari/shitdetector/errors"
	"github.com/kettari/shitdetector/internal/currency"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const cbrfUrl = `https://www.cbr.ru/scripts/XML_daily.asp`

type (
	cbrfProvider struct {
	}
)
func NewCbrfProvider() *cbrfProvider {
	return &cbrfProvider{}
}

func (p cbrfProvider) Fetch(symbol string) (curr *currency.Currency, err error) {
	resp, err := http.Get(cbrfUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	expression := fmt.Sprintf(`<Valute ID="[\w]+"><NumCode>[\d]{3}</NumCode><CharCode>%s</CharCode>.+?<Value>([\d,]+)</Value></Valute>`, symbol)
	r := regexp.MustCompile(expression)
	found := r.FindStringSubmatch(string(raw))
	if len(found) > 1 {
		rateString := strings.Replace(found[1], ",", ".", -1)
		rate, err := strconv.ParseFloat(rateString, 64)
		if err != nil {
			return nil, err
		}
		return &currency.Currency{
			Symbol:  symbol,
			Rate:    rate,
			Created: time.Now().Unix(),
		}, nil
	} else {
		return nil, errors.ErrCurrencyNotFound
	}
}
