package domain

import "github.com/google/uuid"

type TradeItemType int

const (
	ItemTradeType TradeItemType = iota
	RoosterTradeType
)

type TradeItem struct {
	Type    TradeItemType
	Rooster *Rooster // can be nil
	Item    *Item    // can be nil
}

type TradeUser struct {
	ID        ID
	Items     []*TradeItem
	Confirmed bool
}

type Trade struct {
	ID    uuid.UUID
	Users map[ID]*TradeUser
}
