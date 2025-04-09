package trade

import (
	"errors"
	"slices"

	"github.com/acnologla/asuraTrades/internal/core/domain"
)

type TradeUser struct {
	ID        domain.ID
	Items     []*TradeItem
	Confirmed bool
}

func (user *TradeUser) getItemsByType(itemType TradeItemType) []*TradeItem {
	items := []*TradeItem{}
	for _, item := range user.Items {
		if item.Type == itemType {
			items = append(items, item)
		}
	}
	return items
}

func (user *TradeUser) appendItem(item *TradeItem) error {
	if len(user.Items) >= MAXIMUM_TRADE_ITEMS {
		return errors.New("maximum trade items exceeded")
	}

	user.Items = append(user.Items, item)
	return nil
}

func (user *TradeUser) removeItem(item *TradeItem) {
	i := 0
	for ; i < len(user.Items); i++ {
		if user.Items[i] == item {
			break
		}
	}
	user.Items = slices.Delete(user.Items, i, i+1)
}

func (user *TradeUser) addGeneric(item *TradeItem) error {
	for _, it := range user.getItemsByType(item.Type) {
		if item.TradeObject.GetID() == it.TradeObject.GetID() {
			return errors.New("item already added")
		}
	}
	return user.appendItem(item)
}
func (user *TradeUser) addItem(item *TradeItem) error {
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

	itemEntity.Quantity = 1 // We set this quantity to 1 because the user can only add one item at a time

	return user.appendItem(item)
}
