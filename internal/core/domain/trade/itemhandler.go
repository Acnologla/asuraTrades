package trade

import (
	"errors"

	"github.com/google/uuid"
)

type ItemTypeHandler interface {
	Add(user *TradeUser, item *TradeItem) error
	Remove(user *TradeUser, itemID uuid.UUID, tradeItemType TradeItemType) error
}

type RegularItemHandler struct{}

func (h *RegularItemHandler) Add(user *TradeUser, item *TradeItem) error {
	itemEntity := item.Item()

	for _, it := range user.getItemsByType(ItemTradeType) {
		if itemEntity.ID == it.Item().ID {
			if it.Item().Quantity+1 > itemEntity.Quantity {
				return errors.New("item quantity exceeded")
			}
			it.Item().Quantity++
			return nil
		}
	}

	itemEntity.Quantity = 1
	return user.appendItem(item)
}

func (h *RegularItemHandler) Remove(user *TradeUser, itemID uuid.UUID, _ TradeItemType) error {
	for _, item := range user.getItemsByType(ItemTradeType) {
		itemEntity := item.Item()
		if itemEntity.ID == itemID {
			if itemEntity.Quantity > 1 {
				itemEntity.Quantity--
				return nil
			}
			user.removeItem(item)
			return nil
		}
	}
	return errors.New("item not found")
}

type GenericItemHandler struct{}

func (h *GenericItemHandler) Add(user *TradeUser, item *TradeItem) error {
	for _, it := range user.getItemsByType(item.Type) {
		if item.TradeObject.GetID() == it.TradeObject.GetID() {
			return errors.New("item already added")
		}
	}
	return user.appendItem(item)
}

func (h *GenericItemHandler) Remove(user *TradeUser, itemID uuid.UUID, itemType TradeItemType) error {
	for _, item := range user.getItemsByType(itemType) {
		if item.TradeObject.GetID() == itemID {
			user.removeItem(item)
			return nil
		}
	}
	return errors.New("item not found")
}

var itemHandlers = map[TradeItemType]ItemTypeHandler{
	ItemTradeType:    &RegularItemHandler{},
	PetTradeType:     &GenericItemHandler{},
	RoosterTradeType: &GenericItemHandler{},
}
