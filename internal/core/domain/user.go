package domain

import (
	"strconv"

	"github.com/google/uuid"
)

type ID uint64

func (id ID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

func NewID(id string) (ID, error) {
	i, err := strconv.ParseUint(id, 10, 64)
	return ID(i), err
}

type User struct {
	ID ID
	Xp int
}

type UserProfile struct {
	*User
	Roosters []*Rooster
	Items    []*Item
	Pets     []*Pet
}

func NewUserProfile(user *User, roosters []*Rooster, items []*Item, pets []*Pet) *UserProfile {
	return &UserProfile{
		User:     user,
		Roosters: roosters,
		Items:    items,
		Pets:     pets,
	}
}

type UserTrade struct {
	AuthorID ID
	OtherID  ID
	TradeID  uuid.UUID
}

func NewUserTrade(authorID, otherID, tradeID string) (*UserTrade, error) {
	author, err := NewID(authorID)
	if err != nil {
		return nil, err
	}
	other, err := NewID(otherID)
	if err != nil {
		return nil, err
	}

	parsedID, err := uuid.Parse(tradeID)

	if err != nil {
		return nil, err
	}

	return &UserTrade{
		AuthorID: author,
		OtherID:  other,
		TradeID:  parsedID,
	}, nil
}
