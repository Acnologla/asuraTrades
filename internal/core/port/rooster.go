package port

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
)

type RoosterRepository interface {
	Get(ctx context.Context, id domain.ID) (*domain.Rooster, error)
	Delete(ctx context.Context, id domain.ID) error
	Create(ctx context.Context, rooster *domain.Rooster) error
}
