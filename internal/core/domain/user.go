package domain

type ID = uint64

type User struct {
	ID ID
	Xp int
}

type UserProfile struct {
	*User
	Roosters []*Rooster
	Items    []*Item
}
