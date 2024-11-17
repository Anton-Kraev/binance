package domain

import "sync"

type TradeStream struct {
	trades  []Trade
	lastIdx int
	mu      sync.Mutex
}

func NewTradeStream() *TradeStream {
	return &TradeStream{}
}

func (l *TradeStream) GetNew() []Trade {
	l.mu.Lock()
	defer l.mu.Unlock()

	newTrades := l.trades[l.lastIdx:]
	l.lastIdx = len(l.trades)

	return newTrades
}

func (l *TradeStream) GetAll() []Trade {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.trades
}

func (l *TradeStream) Add(trade Trade) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.trades = append(l.trades, trade)
}
