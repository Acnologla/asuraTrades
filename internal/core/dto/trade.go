package dto

import (
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/google/uuid"
)

type TradeItemDTO struct {
	Type   domain.TradeItemType
	ID     uuid.UUID
	User   domain.ID
	Remove bool
}

type UpdateUserStatusDTO struct {
	ID        uuid.UUID
	Confirmed bool
	UserID    domain.ID
}
