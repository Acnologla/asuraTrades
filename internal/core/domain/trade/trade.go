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

func (t *Trade) getUserAndHandler(userID domain.ID, itemType TradeItemType) (*TradeUser, ItemTypeHandler, error) {
	user, err := t.FindUser(userID)
	if err != nil {
		return nil, nil, err
	}

	handler, ok := itemHandlers[itemType]
	if !ok {
		return nil, nil, errors.New("item type not supported")
	}

	return user, handler, nil
}

func (t *Trade) AddItem(userID domain.ID, item *TradeItem) error {
	if !item.TradeObject.IsTradeable() {
		return errors.New("item is not tradeable")
	}

	user, handler, err := t.getUserAndHandler(userID, item.Type)
	if err != nil {
		return err
	}

	return handler.Add(user, item)
}

func (t *Trade) RemoveItem(userID domain.ID, itemID uuid.UUID, itemType TradeItemType) error {

	user, handler, err := t.getUserAndHandler(userID, itemType)
	if err != nil {
		return err
	}

	return handler.Remove(user, itemID, itemType)
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
