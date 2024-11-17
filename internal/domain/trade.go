package domain

import (
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

type Trade struct {
	Time     int64
	Quantity float64
}

func FromBinanceTrade(trade *binance.Trade) Trade {
	if trade == nil {
		return Trade{}
	}

	qty, err := strconv.ParseFloat(trade.Quantity, 64)
	if err != nil {
		qty = 0
	}

	return Trade{
		Time:     trade.Time,
		Quantity: qty,
	}
}

func FromBinanceTradeEvent(event *binance.WsTradeEvent) Trade {
	if event == nil {
		return Trade{}
	}

	qty, err := strconv.ParseFloat(event.Quantity, 64)
	if err != nil {
		qty = 0
	}

	price, err := strconv.ParseFloat(event.Price, 64)
	if err != nil {
		price = 0
	}

	return Trade{
		Time:     event.TradeTime,
		Quantity: price * qty,
	}
}

func TradeFields() string {
	return "Time,Quantity"
}

func (t Trade) String() string {
	return fmt.Sprintf("%s,%.0f",
		time.UnixMilli(t.Time).Format("02.01.2006-15:04:05"),
		t.Quantity,
	)
}

func (t Trade) IsSameTrade(nextTrade Trade) bool {
	return t.Time == nextTrade.Time
}

func (t Trade) Merge(nextTrade Trade) Trade {
	return Trade{
		Time:     t.Time,
		Quantity: t.Quantity + nextTrade.Quantity,
	}
}

func (t Trade) MatchesFilter(qtyThreshold float64) bool {
	isValid := t.Time > 0
	isOverThreshold := t.Quantity > qtyThreshold

	return isValid && isOverThreshold
}
