package asset

type Service interface {
	Get(ticker string) (*Stock, error)
}

type Stock struct {
	ID        string
	Ticker    string
	MarketCap float64
	EPS       float64
	ROE       float64
	Leverage  float64
	EPSRate   float64
}
