package port

import "github.com/acnologla/asuraTrades/internal/core/domain"

type TokenProvider interface {
	GenerateToken(authorID, offerID domain.ID, minutesToExpire int) (string, error)
	ValidateToken(token string) (*domain.UserTrade, error)
}
