package domain

import "github.com/google/uuid"

type Pet struct {
	ID     uuid.UUID
	UserID ID
	Type   int
	Level  int
}

func (r *Pet) GetID() uuid.UUID {
	return r.ID
}

func (r *Pet) IsTradeable() bool {
	return true
}

func NewPet(userID ID, petType, level int) *Pet {
	return &Pet{
		UserID: userID,
		Type:   petType,
		Level:  level,
	}
}
