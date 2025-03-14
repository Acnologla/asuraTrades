package domain

import "github.com/google/uuid"

type Rooster struct {
	ID     uuid.UUID
	UserID ID
	Origin string
	Type   int
	Equip  bool
}

func (r *Rooster) IsTradeable() bool {
	return !r.Equip
}

func GetTradableRoosters(roosters []*Rooster) []*Rooster {
	tradableRoosters := make([]*Rooster, 0, len(roosters))
	for _, rooster := range roosters {
		if rooster.IsTradeable() {
			tradableRoosters = append(tradableRoosters, rooster)
		}
	}
	return tradableRoosters
}

func NewRooster(userID ID, t int, origin string) *Rooster {
	return &Rooster{
		UserID: userID,
		Type:   t,
		Origin: origin,
	}
}
