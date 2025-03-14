package response

import "github.com/acnologla/asuraTrades/internal/core/domain"

type TradeItemResponse struct {
	Type    int                  `json:"type"`
	Rooster *UserRoosterResponse `json:"rooster"`
	Item    *UserItemResponse    `json:"item"`
}

type TradeUserResponse struct {
	ID        string               `json:"id"`
	Confirmed bool                 `json:"confirmed"`
	Items     []*TradeItemResponse `json:"items"`
}

type TradeResponse struct {
	ID    string                        `json:"id"`
	Users map[string]*TradeUserResponse `json:"users"`
}

func NewTradeResponse(trade domain.Trade) *TradeResponse {
	users := make(map[string]*TradeUserResponse, len(trade.Users))
	for id, user := range trade.Users {
		items := make([]*TradeItemResponse, len(user.Items))
		for i, item := range user.Items {
			items[i] = &TradeItemResponse{
				Type: int(item.Type),
			}
			if item.Type == domain.RoosterTradeType {
				items[i].Rooster = &UserRoosterResponse{
					ID:     item.Rooster.ID.String(),
					UserID: user.ID.String(),
					Origin: item.Rooster.Origin,
					Type:   item.Rooster.Type,
				}
			} else {
				items[i].Item = &UserItemResponse{
					ID:       item.Item.ID.String(),
					UserID:   user.ID.String(),
					Quantity: item.Item.Quantity,
					ItemID:   item.Item.ItemID,
					Type:     int(item.Item.Type),
				}
			}
		}
		users[id.String()] = &TradeUserResponse{
			ID:        user.ID.String(),
			Confirmed: user.Confirmed,
			Items:     items,
		}
	}
	return &TradeResponse{
		ID:    trade.ID.String(),
		Users: users,
	}
}
