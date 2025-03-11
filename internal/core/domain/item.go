package domain

import "github.com/google/uuid"

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

type Item struct {
	ID       uuid.UUID
	UserID   ID
	Quantity int
	ItemID   int
	Type     ItemType
}
