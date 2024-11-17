package domain

import (
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"

	"binance/internal/helpers"
)

type Trade struct {
	IsBuyerMaker bool
	Time         int64
	Quantity     float64
	Price        float64
}

// TODO: fix included fields in return, before use
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
		IsBuyerMaker: event.IsBuyerMaker,
		Time:         event.TradeTime,
		Quantity:     price * qty,
		Price:        price,
	}
}

func TradeFields() string {
	return "Trend,Time,Quantity,Price"
}

func (t Trade) String() string {
	// TODO: choose number of digits after '.' accordingly to field size
	return fmt.Sprintf("%s,%s,%.0f,%.0f",
		helpers.GetTrendColor(t.IsBuyerMaker),
		time.UnixMilli(t.Time).Format("02.01.2006-15:04:05"),
		t.Quantity,
		t.Price,
	)
}

func (t Trade) IsSameTrade(nextTrade Trade) bool {
	return t.Time == nextTrade.Time && t.IsBuyerMaker == nextTrade.IsBuyerMaker
}

func (t Trade) Merge(nextTrade Trade) Trade {
	return Trade{
		IsBuyerMaker: t.IsBuyerMaker,
		Time:         t.Time,
		Quantity:     t.Quantity + nextTrade.Quantity,
		Price:        (t.Price + nextTrade.Price) / 2,
	}
}

func (t Trade) MatchesFilter(qtyThreshold float64) bool {
	isValid := t.Time > 0
	isOverThreshold := t.Quantity > qtyThreshold

	return isValid && isOverThreshold
}
