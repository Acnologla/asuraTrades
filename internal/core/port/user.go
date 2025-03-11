package port

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
)

type UserRepository interface {
	Get(ctx context.Context, id domain.ID) (*domain.User, error)
	LockUpdate(ctx context.Context, id domain.ID) (func(err error) error, error)
}
