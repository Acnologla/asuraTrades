package port

import (
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/google/uuid"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

type TradeCache interface {
	Get(id uuid.UUID) (*domain.Trade, error)
	Set(id uuid.UUID, trade *domain.Trade) error
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, trade *domain.Trade) error
}
