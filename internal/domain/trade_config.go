package domain

type TradeConfig struct {
	Symbol       string
	QtyThreshold float64
	Merge        bool
}

func NewTradeConfig(symbol string, qtyThreshold float64, merge bool) TradeConfig {
	return TradeConfig{
		Symbol:       symbol,
		QtyThreshold: qtyThreshold,
		Merge:        merge,
	}
}
