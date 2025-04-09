package cache

import (
	"errors"

	"github.com/acnologla/asuraTrades/internal/core/domain/trade"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/google/uuid"
)

type LocalCache struct {
	trades map[uuid.UUID]*trade.Trade
}

func (l *LocalCache) Get(id uuid.UUID) (*trade.Trade, error) {
	trade, ok := l.trades[id]
	if !ok {
		return nil, errors.New("trade not found")
	}
	return trade, nil
}

func (l *LocalCache) Set(id uuid.UUID, trade *trade.Trade) error {
	l.trades[id] = trade
	return nil
}

func (l *LocalCache) Delete(id uuid.UUID) error {
	delete(l.trades, id)
	return nil
}

func (l *LocalCache) Update(id uuid.UUID, trade *trade.Trade) error {
	l.trades[id] = trade
	return nil
}

func NewLocalCache() port.TradeCache {
	return &LocalCache{
		trades: make(map[uuid.UUID]*trade.Trade),
	}
}
