package domain

import (
	"errors"

	"github.com/google/uuid"
)

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

func NewTradeItemRooster(rooster *Rooster) *TradeItem {
	return &TradeItem{
		Type:    RoosterTradeType,
		Rooster: rooster,
	}
}

func NewTradeItemItem(item *Item) *TradeItem {
	return &TradeItem{
		Type: ItemTradeType,
		Item: item,
	}
}

type TradeUser struct {
	ID        ID
	Items     []*TradeItem
	Confirmed bool
}

func (user *TradeUser) addRoster(item *TradeItem) error {
	if item.Rooster.Equip {
		return errors.New("rooster already equipped")
	}
	for _, it := range user.Items {
		if it.Type == RoosterTradeType && it.Rooster.ID == item.Rooster.ID {
			return errors.New("rooster already added")
		}
	}
	user.Items = append(user.Items, item)
	return nil
}

func (user *TradeUser) addItem(item *TradeItem) error {
	if _, ok := tradeableItemTypes[item.Item.Type]; !ok {
		return errors.New("item type not tradeable")
	}

	for _, it := range user.Items {
		if it.Type == ItemTradeType && it.Item.ID == item.Item.ID {
			it.Item.Quantity++
			return nil
		}
	}
	user.Items = append(user.Items, item)
	return nil
}

type Trade struct {
	ID    uuid.UUID
	Users map[ID]*TradeUser
}

func (t *Trade) AddItem(userID ID, item *TradeItem) error {
	user, ok := t.Users[userID]
	if !ok {
		return errors.New("user not found")
	}

	if item.Type == RoosterTradeType {
		return user.addRoster(item)
	}

	return user.addItem(item)
}

func (t *Trade) UpdateUserStatus(userID ID, confirmed bool) error {
	user, ok := t.Users[userID]
	if !ok {
		return errors.New("user not found")
	}

	user.Confirmed = confirmed
	return nil
}

func (t *Trade) Done() bool {
	for _, user := range t.Users {
		if !user.Confirmed {
			return false
		}
	}
	return true
}

func (t *Trade) RemoveItem(userID ID, itemID uuid.UUID) error {
	user, ok := t.Users[userID]
	if !ok {
		return errors.New("user not found")
	}

	for i, item := range user.Items {
		if item.Type == ItemTradeType && item.Item.ID == itemID {
			if item.Item.Quantity > 1 {
				item.Item.Quantity--
				return nil
			}

			user.Items = append(user.Items[:i], user.Items[i+1:]...)
			return nil
		}
	}

	return errors.New("item not found")
}

func NewTrade(id uuid.UUID, author, other ID) *Trade {
	return &Trade{
		ID: id,
		Users: map[ID]*TradeUser{
			author: {
				ID: author,
			},
			other: {
				ID: other,
			},
		},
	}
}
