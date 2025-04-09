package port

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/google/uuid"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

type PetRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Pet, error)
	GetUserPets(ctx context.Context, id domain.ID) ([]*domain.Pet, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Create(ctx context.Context, pet *domain.Pet) error
}
