package domain

import (
	"errors"
	"slices"

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
	rooster *Rooster
	item    *Item
}

func (t *TradeItem) Rooster() *Rooster {
	if t.Type != RoosterTradeType {
		panic("trade item is not a rooster")
	}
	return t.rooster
}

func (t *TradeItem) Item() *Item {
	if t.Type != ItemTradeType {
		panic("trade item is not an item")
	}
	return t.item
}

func NewTradeItemRooster(rooster *Rooster) *TradeItem {
	return &TradeItem{
		Type:    RoosterTradeType,
		rooster: rooster,
	}
}

func NewTradeItemItem(item *Item) *TradeItem {
	return &TradeItem{
		Type: ItemTradeType,
		item: item,
	}
}

type TradeUser struct {
	ID        ID
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

func (user *TradeUser) addRoster(item *TradeItem) error {
	rooster := item.Rooster()
	if !rooster.IsTradeable() {
		return errors.New("rooster is not tradeable")
	}
	for _, it := range user.getItemsByType(RoosterTradeType) {
		if rooster.ID == it.Rooster().ID {
			return errors.New("rooster already added")
		}
	}
	return user.appendItem(item)
}

func (user *TradeUser) addItem(item *TradeItem) error {
	itemEntity := item.Item()
	if !itemEntity.IsTradeable() {
		return errors.New("item is not tradeable")
	}

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

type Trade struct {
	ID    uuid.UUID
	Users map[ID]*TradeUser
}

func (t *Trade) FindUser(userID ID) (*TradeUser, error) {
	user, ok := t.Users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (t *Trade) AddItem(userID ID, item *TradeItem) error {
	user, err := t.FindUser(userID)
	if err != nil {
		return err
	}

	if item.Type == RoosterTradeType {
		return user.addRoster(item)
	}

	return user.addItem(item)
}

func (t *Trade) UpdateUserStatus(userID ID, confirmed bool) error {
	user, err := t.FindUser(userID)
	if err != nil {
		return err
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

func (t *Trade) removeRooster(user *TradeUser, itemID uuid.UUID) error {
	for _, item := range user.getItemsByType(RoosterTradeType) {
		if item.Rooster().ID == itemID {
			user.removeItem(item)
			return nil
		}
	}
	return errors.New("rooster not found")
}

func (t *Trade) RemoveItem(userID ID, itemID uuid.UUID, itemType TradeItemType) error {
	user, err := t.FindUser(userID)
	if err != nil {
		return err
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
