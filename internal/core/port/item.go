package port

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/google/uuid"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

type ItemRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Item, error)
	GetUserItems(ctx context.Context, id domain.ID) ([]*domain.Item, error)
	Remove(ctx context.Context, id uuid.UUID, quantity int) error
	Add(ctx context.Context, item *domain.Item, quantity int) error
}
