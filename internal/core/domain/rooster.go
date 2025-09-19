package domain

import "github.com/google/uuid"

type Rooster struct {
	ID      uuid.UUID
	UserID  ID
	Origin  string
	Type    int
	Equip   bool
	Special bool
}

func (r *Rooster) GetID() uuid.UUID {
	return r.ID
}

func (r *Rooster) IsTradeable() bool {
	return !r.Equip
}

func NewRooster(userID ID, t int, origin string, special bool) *Rooster {
	return &Rooster{
		UserID:  userID,
		Type:    t,
		Origin:  origin,
		Special: special,
	}
}
