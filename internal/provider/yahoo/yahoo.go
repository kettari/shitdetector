package yahoo

import (
	"encoding/json"
	"fmt"
	"github.com/kettari/shitdetector/errors"
	"github.com/kettari/shitdetector/internal/asset"
	"net/http"
	"time"
)

const yahooUrl = `https://query1.finance.yahoo.com/v10/finance/quoteSummary/%s?modules=price,defaultKeyStatistics,financialData,earningsTrend,earningsHistory`

type (
	yahooProvider struct {
	}
	yahooResponse struct {
		Summary QuoteSummary `json:"quoteSummary"`
	}
	QuoteSummary struct {
		Result []*QuoteResult `json:"result"`
	}
	QuoteResult struct {
		DefaultKeyStatistics DefaultKeyStatisticsStruct `json:"defaultKeyStatistics"`
		FinancialData        FinancialDataStruct        `json:"financialData"`
		Price                PriceStruct                `json:"price"`
	}
	DefaultKeyStatisticsStruct struct {
		TrailingEPS QuoteRecord `json:"trailingEps"`
	}
	FinancialDataStruct struct {
		ReturnOnEquity QuoteRecord `json:"returnOnEquity"`
		DebtToEquity   QuoteRecord `json:"debtToEquity"`
	}
	PriceStruct struct {
		MarketCap QuoteRecord `json:"marketCap"`
		ShortName string      `json:"shortName"`
	}
	QuoteRecord struct {
		Raw float64 `json:"raw"`
		Fmt string  `json:"fmt"`
	}
)

func NewYahooProvider() *yahooProvider {
	return &yahooProvider{}
}

func (p yahooProvider) Fetch(ticker string) (stock *asset.Stock, err error) {
	quoteSummary, err := fetchQuoteSummary(fmt.Sprintf(yahooUrl, ticker))
	if err != nil {
		return nil, err
	}
	if len(quoteSummary.Summary.Result) == 0 {
		return nil, errors.ErrQuoteSummaryEmpty
	}

	stock = &asset.Stock{
		Ticker:    ticker,
		ShortName: quoteSummary.Summary.Result[0].Price.ShortName,
		Created:   time.Now().Unix(),
		MarketCap: quoteSummary.Summary.Result[0].Price.MarketCap.Raw,
		EPS:       quoteSummary.Summary.Result[0].DefaultKeyStatistics.TrailingEPS.Raw,
		ROE:       quoteSummary.Summary.Result[0].FinancialData.ReturnOnEquity.Raw,
		Leverage:  quoteSummary.Summary.Result[0].FinancialData.DebtToEquity.Raw,
		EPSRate:   0,
	}

	return stock, nil
}

func fetchQuoteSummary(url string) (*yahooResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	yahooResp := &yahooResponse{}
	if err = json.NewDecoder(resp.Body).Decode(yahooResp); err != nil {
		return nil, err
	}

	return yahooResp, nil
}
