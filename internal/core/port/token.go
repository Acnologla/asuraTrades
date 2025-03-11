package port

import "github.com/acnologla/asuraTrades/internal/core/domain"

type TokenService interface {
	GenerateToken(id domain.ID, minutesToExpire int) (string, error)
	ValidateToken(token string) (domain.ID, error)
}
