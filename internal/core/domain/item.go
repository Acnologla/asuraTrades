package domain

import (
	"github.com/google/uuid"
)

type ItemType int

const (
	_ ItemType = iota
	LootboxType
	NormalType
	CosmeticType
	KeyType
	ShardType
	AchievementType
	SurvivalKeyType
)

var tradeableItemTypes = map[ItemType]struct{}{
	NormalType:   {},
	CosmeticType: {},
	ShardType:    {},
} //use struct instead of bool to save memory

type Item struct {
	ID       uuid.UUID
	UserID   ID
	Quantity int
	ItemID   int
	Type     ItemType
}

func (i *Item) GetID() uuid.UUID {
	return i.ID
}

func (i *Item) IsTradeable() bool {
	_, ok := tradeableItemTypes[i.Type]
	return ok
}

func NewItem(userID ID, itemID int, t ItemType) *Item {
	return &Item{
		UserID: userID,
		ItemID: itemID,
		Type:   t,
	}
}
