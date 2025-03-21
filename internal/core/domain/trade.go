package domain

import (
	"errors"

	"github.com/google/uuid"
)

const MAXIMUM_TRADE_ITEMS = 12

type TradeItemType int

const (
	ItemTradeType TradeItemType = iota
	RoosterTradeType
)

func (t TradeItemType) String() string {
	switch t {
	case ItemTradeType:
		return "item"
	case RoosterTradeType:
		return "rooster"
	}

	return ""
}

type TradeItem struct {
	Type    TradeItemType
	Rooster *Rooster // Will be nil if Type is ItemTradeType
	Item    *Item    // Will be nil if Type is RoosterTradeType
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

func (user *TradeUser) appendItem(item *TradeItem) error {
	if len(user.Items) >= MAXIMUM_TRADE_ITEMS {
		return errors.New("maximum trade items exceeded")
	}

	user.Items = append(user.Items, item)
	return nil
}

func (user *TradeUser) addRoster(item *TradeItem) error {
	if !item.Rooster.IsTradeable() {
		return errors.New("rooster is not tradeable")
	}
	for _, it := range user.Items {
		if it.Type == RoosterTradeType && it.Rooster.ID == item.Rooster.ID {
			return errors.New("rooster already added")
		}
	}
	return user.appendItem(item)
}

func (user *TradeUser) addItem(item *TradeItem) error {
	if !item.Item.IsTradeable() {
		return errors.New("item is not tradeable")
	}

	for _, it := range user.Items {
		if it.Type == ItemTradeType && it.Item.ID == item.Item.ID {
			if it.Item.Quantity+1 > item.Item.Quantity {
				return errors.New("item quantity exceeded")
			}
			it.Item.Quantity++
			return nil
		}
	}

	item.Item.Quantity = 1 // We set this quantity to 1 because the user can only add one item at a time

	return user.appendItem(item)
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

func (t *Trade) removeItem(user *TradeUser, itemID uuid.UUID) error {
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

func (t *Trade) removeRooster(user *TradeUser, itemID uuid.UUID) error {
	for i, item := range user.Items {
		if item.Type == RoosterTradeType && item.Rooster.ID == itemID {
			user.Items = append(user.Items[:i], user.Items[i+1:]...)
			return nil
		}
	}
	return errors.New("rooster not found")
}

func (t *Trade) RemoveItem(userID ID, itemID uuid.UUID, itemType TradeItemType) error {
	user, ok := t.Users[userID]

	if !ok {
		return errors.New("user not found")
	}

	if itemType == RoosterTradeType {
		return t.removeRooster(user, itemID)
	}
	return t.removeItem(user, itemID)
}

func NewTrade(id uuid.UUID, author, other ID) *Trade {
	return &Trade{
		ID: id,
		Users: map[ID]*TradeUser{
			author: {
				ID:    author,
				Items: []*TradeItem{},
			},
			other: {
				ID:    other,
				Items: []*TradeItem{},
			},
		},
	}
}
