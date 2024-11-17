package domain

type TradeConfig struct {
	Symbol       string
	QtyThreshold float64
}

func NewTradeConfig(symbol string, qtyThreshold float64) TradeConfig {
	return TradeConfig{
		Symbol:       symbol,
		QtyThreshold: qtyThreshold,
	}
}
