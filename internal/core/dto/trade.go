package dto

import (
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/domain/trade"
	"github.com/google/uuid"
)

type TradeItemDTO struct {
	Type   trade.TradeItemType
	ID     uuid.UUID
	ItemID uuid.UUID
	User   domain.ID
	Remove bool
}

type UpdateUserStatusDTO struct {
	ID        uuid.UUID
	Confirmed bool
	User      domain.ID
}

func NewTradeItemDTO(t int, ID, itemID uuid.UUID, user domain.ID, remove bool) *TradeItemDTO {
	return &TradeItemDTO{
		Type:   trade.TradeItemType(t),
		ID:     ID,
		ItemID: itemID,
		User:   user,
		Remove: remove,
	}
}

func NewUpdateUserStatusDTO(ID uuid.UUID, confirmed bool, user domain.ID) *UpdateUserStatusDTO {
	return &UpdateUserStatusDTO{
		ID:        ID,
		Confirmed: confirmed,
		User:      user,
	}
}
