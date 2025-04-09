package response

import (
	"github.com/acnologla/asuraTrades/internal/core/domain"
)

type UserItemResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Quantity int    `json:"quantity"`
	ItemID   int    `json:"item_id"`
	Type     int    `json:"type"`
}

type UserRoosterResponse struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Origin string `json:"origin"`
	Type   int    `json:"type"`
}

type UserPetResponse struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Type   int    `json:"type"`
}

type UserTokenResponse struct {
	ID       string                 `json:"id"`
	OtherID  string                 `json:"other_id"`
	Xp       int                    `json:"xp"`
	Roosters []*UserRoosterResponse `json:"roosters"`
	Items    []*UserItemResponse    `json:"items"`
	Pets     []*UserPetResponse     `json:"pets"`
}

func NewUserTokenResponse(userTrade *domain.UserTrade, userProfile *domain.UserProfile) *UserTokenResponse {
	roosters := make([]*UserRoosterResponse, len(userProfile.Roosters))
	uID := userProfile.ID.String()
	for i, rooster := range userProfile.Roosters {
		roosters[i] = &UserRoosterResponse{
			ID:     rooster.ID.String(),
			UserID: uID,
			Origin: rooster.Origin,
			Type:   rooster.Type,
		}
	}
	items := make([]*UserItemResponse, len(userProfile.Items))
	for i, item := range userProfile.Items {
		items[i] = &UserItemResponse{
			ID:       item.ID.String(),
			UserID:   uID,
			Quantity: item.Quantity,
			ItemID:   item.ItemID,
			Type:     int(item.Type),
		}
	}
	pets := make([]*UserPetResponse, len(userProfile.Pets))
	for i, pet := range userProfile.Pets {
		pets[i] = &UserPetResponse{
			ID:     pet.ID.String(),
			UserID: uID,
			Type:   int(pet.Type),
		}
	}
	return &UserTokenResponse{
		ID:       userProfile.ID.String(),
		OtherID:  userTrade.OtherID.String(),
		Xp:       userProfile.Xp,
		Roosters: roosters,
		Items:    items,
		Pets:     pets,
	}
}
