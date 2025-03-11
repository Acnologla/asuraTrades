package response

import (
	"github.com/acnologla/asuraTrades/internal/core/domain"
)

type ItemUserResponse struct {
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

type UserTokenResponse struct {
	ID       string                 `json:"id"`
	OtherID  string                 `json:"other_id"`
	Xp       int                    `json:"xp"`
	Roosters []*UserRoosterResponse `json:"roosters"`
	Items    []*ItemUserResponse    `json:"items"`
}

func NewUserTokenResponse(userTrade *domain.UserTrade, userProfile *domain.UserProfile) *UserTokenResponse {
	roosters := make([]*UserRoosterResponse, len(userProfile.Roosters))
	for i, rooster := range userProfile.Roosters {
		roosters[i] = &UserRoosterResponse{
			ID:     rooster.ID.String(),
			UserID: userProfile.ID.String(),
			Origin: rooster.Origin,
			Type:   rooster.Type,
		}
	}
	items := make([]*ItemUserResponse, len(userProfile.Items))
	for i, item := range userProfile.Items {
		items[i] = &ItemUserResponse{
			ID:       item.ID.String(),
			UserID:   userProfile.ID.String(),
			Quantity: item.Quantity,
			ItemID:   item.ItemID,
			Type:     int(item.Type),
		}
	}
	return &UserTokenResponse{
		ID:       userProfile.ID.String(),
		OtherID:  userTrade.OtherID.String(),
		Xp:       userProfile.Xp,
		Roosters: roosters,
		Items:    items,
	}
}
