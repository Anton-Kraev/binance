package main

import (
	"fmt"
	"log"
	"os"

	"github.com/adshao/go-binance/v2"

	"binance/internal/domain"
	"binance/internal/market"
	"binance/internal/view"
)

var (
	symbol       = "BTCUSDT"
	qtyThreshold = 100000.
)

func main() {
	binanceClient := binance.NewClient("", "")
	tradeStream := domain.NewTradeStream()
	tradeConfig := domain.NewTradeConfig(symbol, qtyThreshold)

	tradeService := market.NewTradesService(binanceClient, tradeStream, tradeConfig)

	done, err := tradeService.StartReceiving()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s>%.0f\n", symbol, qtyThreshold)
	if err = view.RunWritingTrades(done, os.Stdout, tradeStream); err != nil {
		log.Fatalln(err)
	}
}
