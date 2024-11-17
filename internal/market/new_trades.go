package market

import (
	"log"

	"github.com/adshao/go-binance/v2"

	"binance/internal/domain"
)

type tradeStream interface {
	Add(trade domain.Trade)
}

type TradesService struct {
	client *binance.Client
	trades tradeStream
	config domain.TradeConfig
}

func NewTradesService(client *binance.Client, trades tradeStream, config domain.TradeConfig) *TradesService {
	return &TradesService{
		client: client,
		trades: trades,
		config: config,
	}
}

func (c *TradesService) StartReceiving() (done <-chan struct{}, err error) {
	var prevTrade domain.Trade

	wsTradeHandler := func(event *binance.WsTradeEvent) {
		currTrade := domain.FromBinanceTradeEvent(event)

		if prevTrade.IsSameTrade(currTrade) {
			prevTrade = prevTrade.Merge(currTrade)

			return
		}

		if prevTrade.MatchesFilter(c.config.QtyThreshold) {
			c.trades.Add(prevTrade)
		}

		prevTrade = currTrade
	}

	errHandler := func(err error) {
		log.Println(err)
	}

	done, _, err = binance.WsTradeServe(c.config.Symbol, wsTradeHandler, errHandler)

	return
}
