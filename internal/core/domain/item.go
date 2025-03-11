package domain

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
	ID       uint64
	UserID   uint64
	Quantity int
	ItemID   int
	Type     ItemType
}
