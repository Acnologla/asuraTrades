package dto

import (
	"github.com/acnologla/asuraTrades/internal/core/domain"
)

type ItemUserDto struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Quantity int    `json:"quantity"`
	ItemID   int    `json:"item_id"`
	Type     int    `json:"type"`
}

type UserRoosterDto struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Origin string `json:"origin"`
	Type   int    `json:"type"`
}

type UserProfileDTO struct {
	ID       string            `json:"id"`
	Xp       int               `json:"xp"`
	Roosters []*UserRoosterDto `json:"roosters"`
	Items    []*ItemUserDto    `json:"items"`
}

func NewUserProfileDTO(userProfile *domain.UserProfile) *UserProfileDTO {
	roosters := make([]*UserRoosterDto, len(userProfile.Roosters))
	for i, rooster := range userProfile.Roosters {
		roosters[i] = &UserRoosterDto{
			ID:     rooster.ID.String(),
			UserID: userProfile.ID.String(),
			Origin: rooster.Origin,
			Type:   rooster.Type,
		}
	}
	items := make([]*ItemUserDto, len(userProfile.Items))
	for i, item := range userProfile.Items {
		items[i] = &ItemUserDto{
			ID:       item.ID.String(),
			UserID:   userProfile.ID.String(),
			Quantity: item.Quantity,
			ItemID:   item.ItemID,
			Type:     int(item.Type),
		}
	}
	return &UserProfileDTO{
		ID:       userProfile.ID.String(),
		Xp:       userProfile.Xp,
		Roosters: roosters,
		Items:    items,
	}
}
