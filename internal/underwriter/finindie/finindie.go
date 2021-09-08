package finindie

import (
	"github.com/kettari/shitdetector/internal/asset"
	"github.com/kettari/shitdetector/internal/underwriter"
)

const billion = 1000000000

type finindieUnderwriter struct {
	stock *asset.Stock
}

func NewFinindieUnderwriter(stock *asset.Stock) *finindieUnderwriter {
	return &finindieUnderwriter{stock: stock}
}

func (u finindieUnderwriter) Score() (score *underwriter.Score) {
	score = &underwriter.Score{Stock: u.stock}

	// MarketCap
	marketCapCriteria := underwriter.ScoreCriteria{
		Name:  "Рыночная капитализация",
		Value: 1,
	}
	if u.stock.MarketCap >= 30*billion {
		marketCapCriteria.Value = 5
	} else if u.stock.MarketCap >= 20*billion {
		marketCapCriteria.Value = 4
	} else if u.stock.MarketCap >= 10*billion {
		marketCapCriteria.Value = 3
	} else if u.stock.MarketCap >= 2*billion {
		marketCapCriteria.Value = 2
	}
	score.Criterias = append(score.Criterias, &marketCapCriteria)

	// Profitability
	profitabilityCriteria := underwriter.ScoreCriteria{
		Name:  "Прибыльность бзнеса (EPS ttm)",
		Value: 1,
	}
	if u.stock.EPS > 0 {
		profitabilityCriteria.Value = 5
	}
	score.Criterias = append(score.Criterias, &profitabilityCriteria)

	// Return on Equity
	ROECriteria := underwriter.ScoreCriteria{
		Name:  "Рентабельность капитала (ROE)",
		Value: 1,
	}
	if u.stock.ROE >= 0.2 {
		ROECriteria.Value = 5
	} else if u.stock.ROE >= 0.1 {
		ROECriteria.Value = 4
	} else if u.stock.ROE >= 0.05 {
		ROECriteria.Value = 3
	} else if u.stock.ROE >= 0.0 {
		ROECriteria.Value = 2
	}
	score.Criterias = append(score.Criterias, &ROECriteria)

	// Leverage
	LeverageCriteria := underwriter.ScoreCriteria{
		Name:  "Леверидж (Debt/Equity)",
		Value: 5,
	}
	if u.stock.Leverage >= 300 {
		LeverageCriteria.Value = 1
	} else if u.stock.Leverage >= 110 {
		LeverageCriteria.Value = 2
	} else if u.stock.Leverage >= 80 {
		LeverageCriteria.Value = 3
	} else if u.stock.Leverage >= 25 {
		LeverageCriteria.Value = 4
	}
	score.Criterias = append(score.Criterias, &LeverageCriteria)

	return score
}
