package domain

import (
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

type Trade struct {
	ID       int64
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
		ID:       trade.ID,
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
		ID:       event.TradeID,
		Time:     event.TradeTime,
		Quantity: price * qty,
	}
}

func TradeFields() string {
	return "ID,Time,Quantity"
}

func (t Trade) String() string {
	return fmt.Sprintf("%d,%s,%.0f",
		t.ID,
		time.UnixMilli(t.Time).Format("02.01.2006-15:04:05"),
		t.Quantity,
	)
}

func (t Trade) IsSameTrade(nextTrade Trade) bool {
	return t.Time == nextTrade.Time
}

func (t Trade) Merge(nextTrade Trade) Trade {
	return Trade{
		ID:       t.ID,
		Time:     t.Time,
		Quantity: t.Quantity + nextTrade.Quantity,
	}
}

func (t Trade) MatchesFilter(qtyThreshold float64) bool {
	isValid := t.ID > 0 && t.Time > 0
	isOverThreshold := t.Quantity > qtyThreshold

	return isValid && isOverThreshold
}
