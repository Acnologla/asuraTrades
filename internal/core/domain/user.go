package domain

import "strconv"

type ID uint64

func (id ID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

type User struct {
	ID ID
	Xp int
}

type UserProfile struct {
	*User
	Roosters []*Rooster
	Items    []*Item
}
