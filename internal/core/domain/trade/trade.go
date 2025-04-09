package trade

import (
	"errors"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/google/uuid"
)

const MAXIMUM_TRADE_ITEMS = 12

type Trade struct {
	ID    uuid.UUID
	Users map[domain.ID]*TradeUser
}

func (t *Trade) FindUser(userID domain.ID) (*TradeUser, error) {
	user, ok := t.Users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (t *Trade) AddItem(userID domain.ID, item *TradeItem) error {
	if !item.TradeObject.IsTradeable() {
		return errors.New("item is not tradeable")
	}

	user, err := t.FindUser(userID)
	if err != nil {
		return err
	}

	if item.Type == ItemTradeType {
		return user.addItem(item)
	}

	return user.addGeneric(item)
}

func (t *Trade) UpdateUserStatus(userID domain.ID, confirmed bool) error {
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

func (t *Trade) removeGeneric(user *TradeUser, itemID uuid.UUID, itemType TradeItemType) error {
	for _, item := range user.getItemsByType(itemType) {
		if item.TradeObject.GetID() == itemID {
			user.removeItem(item)
			return nil
		}
	}
	return errors.New("item not found")
}

func (t *Trade) RemoveItem(userID domain.ID, itemID uuid.UUID, itemType TradeItemType) error {
	user, err := t.FindUser(userID)
	if err != nil {
		return err
	}

	if itemType == ItemTradeType {
		return t.removeItem(user, itemID)
	}

	return t.removeGeneric(user, itemID, itemType)
}

func NewTrade(id uuid.UUID, author, other domain.ID) *Trade {
	return &Trade{
		ID: id,
		Users: map[domain.ID]*TradeUser{
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
