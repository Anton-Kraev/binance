package view

import (
	"io"

	"binance/internal/domain"
)

type tradeStream interface {
	GetNew() []domain.Trade
}

func RunWritingTrades(done <-chan struct{}, writer io.Writer, trades tradeStream) error {
	_, err := writer.Write([]byte(domain.TradeFields() + "\n"))
	if err != nil {
		return err
	}

	for {
		select {
		case <-done:
			return nil
		default:
			for _, trade := range trades.GetNew() {
				_, err = writer.Write([]byte(trade.String() + "\n"))
				if err != nil {
					return err
				}
			}
		}
	}
}
