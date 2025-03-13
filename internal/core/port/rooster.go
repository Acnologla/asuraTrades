package port

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/google/uuid"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

type RoosterRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Rooster, error)
	GetUserRoosters(ctx context.Context, id domain.ID) ([]*domain.Rooster, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Create(ctx context.Context, rooster *domain.Rooster) error
}
