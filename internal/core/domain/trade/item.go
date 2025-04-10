package trade

import (
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/google/uuid"
)

type TradeItemType int

const (
	ItemTradeType TradeItemType = iota
	RoosterTradeType
	PetTradeType
)

func (t TradeItemType) String() string {
	return map[TradeItemType]string{
		ItemTradeType:    "item",
		RoosterTradeType: "rooster",
		PetTradeType:     "pet",
	}[t]
}

type Tradeable interface {
	IsTradeable() bool
	GetID() uuid.UUID
}

func GetTradableEntities[T Tradeable](entities []T) []T {
	tradable := make([]T, 0, len(entities))
	for _, entity := range entities {
		if entity.IsTradeable() {
			tradable = append(tradable, entity)
		}
	}
	return tradable
}

type TradeItem struct {
	Type        TradeItemType
	TradeObject Tradeable
}

func (t *TradeItem) Rooster() *domain.Rooster {
	if t.Type != RoosterTradeType {
		panic("trade item is not a rooster")
	}
	return t.TradeObject.(*domain.Rooster)
}

func (t *TradeItem) Pet() *domain.Pet {
	if t.Type != PetTradeType {
		panic("trade item is not a pet")
	}
	return t.TradeObject.(*domain.Pet)
}

func (t *TradeItem) Item() *domain.Item {
	if t.Type != ItemTradeType {
		panic("trade item is not an item")
	}
	return t.TradeObject.(*domain.Item)
}

func NewTradeItemPet(pet *domain.Pet) *TradeItem {
	return &TradeItem{
		Type:        PetTradeType,
		TradeObject: pet,
	}
}

func NewTradeItemRooster(rooster *domain.Rooster) *TradeItem {
	return &TradeItem{
		Type:        RoosterTradeType,
		TradeObject: rooster,
	}
}

func NewTradeItemItem(item *domain.Item) *TradeItem {
	return &TradeItem{
		Type:        ItemTradeType,
		TradeObject: item,
	}
}
