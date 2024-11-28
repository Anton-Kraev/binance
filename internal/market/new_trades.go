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
	var tradeBuffer []domain.Trade

	wsTradeHandler := func(event *binance.WsAggTradeEvent) {
		currTrade := domain.FromBinanceAggTradeEvent(event)

		if tradeBuffer == nil || tradeBuffer[0].IsSameTrade(currTrade) && c.config.Merge {
			tradeBuffer = append(tradeBuffer, currTrade)

			return
		}

		mergedTrade := tradeBuffer[0]

		for _, trade := range tradeBuffer[1:] {
			mergedTrade = mergedTrade.Merge(trade)
		}

		tradeBuffer = nil
		if !c.config.Merge {
			tradeBuffer = append(tradeBuffer, currTrade)
		}

		if mergedTrade.MatchesFilter(c.config.QtyThreshold) {
			c.trades.Add(mergedTrade)
		}
	}

	errHandler := func(err error) {
		log.Println(err)
	}

	done, _, err = binance.WsAggTradeServe(c.config.Symbol, wsTradeHandler, errHandler)

	return
}
