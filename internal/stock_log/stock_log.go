package stock_log

type (
	StockLog struct {
		ID      string
		Created int64
		Ticker  string
	}
	Service interface {
		Log(string) error
		Cleanup() error
		Stats() (TickerStats, error)
	}
	TickerStats []*TickerStat
	TickerStat  struct {
		Ticker string
		Count  int64
	}
)

func (t TickerStats) Len() int           { return len(t) }
func (t TickerStats) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TickerStats) Less(i, j int) bool { return t[i].Count > t[j].Count }
