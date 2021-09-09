package yahoo

import (
	"encoding/json"
	"fmt"
	"github.com/kettari/shitdetector/errors"
	"github.com/kettari/shitdetector/internal/asset"
	"net/http"
	"time"
)

const yahooUrl = `https://query1.finance.yahoo.com/v10/finance/quoteSummary/%s?modules=price,defaultKeyStatistics,financialData,earningsTrend`

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
		EarningsTrend        EarningsTrendStruct        `json:"earningsTrend"`
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
		Currency  string      `json:"currency"`
	}
	EarningsTrendStruct struct {
		Trend []*TrendStruct `json:"trend"`
	}
	TrendStruct struct {
		Period string      `json:"period"`
		Growth QuoteRecord `json:"growth"`
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

	// Calculate EPS rate
	pastEPS := float64(0.0)
	futureEPS := float64(0.0)
	for _, trend := range quoteSummary.Summary.Result[0].EarningsTrend.Trend {
		if trend.Period == "-5y" {
			pastEPS = trend.Growth.Raw
		}
		if trend.Period == "+5y" {
			futureEPS = trend.Growth.Raw
		}
	}
	epsRate := (pastEPS + futureEPS) / 2

	stock = &asset.Stock{
		Ticker:    ticker,
		ShortName: quoteSummary.Summary.Result[0].Price.ShortName,
		Created:   time.Now().Unix(),
		MarketCap: quoteSummary.Summary.Result[0].Price.MarketCap.Raw,
		EPS:       quoteSummary.Summary.Result[0].DefaultKeyStatistics.TrailingEPS.Raw,
		ROE:       quoteSummary.Summary.Result[0].FinancialData.ReturnOnEquity.Raw,
		Leverage:  quoteSummary.Summary.Result[0].FinancialData.DebtToEquity.Raw,
		EPSRate:   epsRate,
		Currency:  quoteSummary.Summary.Result[0].Price.Currency,
	}
	if stock.Currency == "" {
		stock.Currency = "USD"
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
