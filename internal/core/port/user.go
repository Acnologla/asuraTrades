package port

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

type UserRepository interface {
	Get(ctx context.Context, id domain.ID) (*domain.User, error)
}

type UserTradeTxAdapters struct {
	UserRepository    UserRepository
	ItemRepository    ItemRepository
	RoosterRepository RoosterRepository
	PetRepository     PetRepository
}

type TradeTxProvider interface {
	Transact(ctx context.Context, txFunc func(adapters UserTradeTxAdapters, lock func(domain.ID) error) error) error
}
