package port

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

type ItemRepository interface {
	Get(ctx context.Context, id domain.ID) (*domain.Item, error)
	GetUserItems(ctx context.Context, id domain.ID) ([]*domain.Item, error)
	Remove(ctx context.Context, id domain.ID) error
	Add(ctx context.Context, item *domain.Item) error
}
