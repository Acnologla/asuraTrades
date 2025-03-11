package domain

import "github.com/google/uuid"

type Rooster struct {
	ID     uuid.UUID
	UserID ID
	Origin string
	Type   int
}
