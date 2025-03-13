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

func GetTradableItems(items []*Item) []*Item {
	tradableItems := make([]*Item, 0, len(items))
	for _, item := range items {
		if _, ok := tradeableItemTypes[item.Type]; ok {
			tradableItems = append(tradableItems, item)
		}
	}
	return tradableItems
}
