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
