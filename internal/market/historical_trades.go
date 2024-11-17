package market

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2"

	"binance/internal/domain"
)

type Client struct {
	client *binance.Client
	trades domain.TradeStream
}

func NewClient(client *binance.Client, trades domain.TradeStream) *Client {
	return &Client{client: client, trades: trades}
}

func (c *Client) GetDayFirstTradeID() (int64, error) {
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	binanceDayStart := binance.FormatTimestamp(dayStart)

	dayFirstTrade, err := c.client.NewAggTradesService().
		Symbol("BTCUSDT").
		StartTime(binanceDayStart).
		Limit(1).
		Do(context.TODO())
	if err != nil {
		return 0, err
	}

	return dayFirstTrade[0].FirstTradeID, nil
}

func (c *Client) GetDayTrades(toID int64) ([]domain.Trade, error) {
	firstDayTradeID, err := c.GetDayFirstTradeID()
	if err != nil {
		return nil, err
	}

	var (
		dayTrades     []domain.Trade
		ctx           = context.TODO()
		tradesService = c.client.NewHistoricalTradesService().Symbol("BTCUSDT")
	)

	for ID := firstDayTradeID; ID < toID; ID += 1000 {
		tradesBatch, err := tradesService.FromID(ID).Limit(min(1000, int(toID-ID))).Do(ctx)
		if err != nil {
			return nil, err
		}

		for _, trade := range tradesBatch {
			dayTrades = append(dayTrades, domain.FromBinanceTrade(trade))
		}

		fmt.Println(ID, time.Unix(0, dayTrades[len(dayTrades)-1].Time*1_000_000))
	}

	return dayTrades, nil
}
