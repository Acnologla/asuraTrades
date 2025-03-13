package domain

import "github.com/google/uuid"

type Rooster struct {
	ID     uuid.UUID
	UserID ID
	Origin string
	Type   int
	Equip  bool
}

func GetTradableRoosters(roosters []*Rooster) []*Rooster {
	tradableRoosters := make([]*Rooster, 0, len(roosters))
	for _, rooster := range roosters {
		if !rooster.Equip {
			tradableRoosters = append(tradableRoosters, rooster)
		}
	}
	return tradableRoosters
}
