package port

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

type UserRepository interface {
	Get(ctx context.Context, id domain.ID) (*domain.User, error)
	LockUpdate(ctx context.Context, userID domain.ID) (func(err error) error, error)
}
