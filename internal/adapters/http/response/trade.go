package response

import (
	tradedomain "github.com/acnologla/asuraTrades/internal/core/domain/trade"
	"github.com/google/uuid"
)

type TradeItemResponse struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type TradeUserResponse struct {
	ID        string               `json:"id"`
	Confirmed bool                 `json:"confirmed"`
	Items     []*TradeItemResponse `json:"items"`
}

type TradeResponse struct {
	Type  string                        `json:"type"` //maybe changer this type later
	ID    string                        `json:"id"`
	Users map[string]*TradeUserResponse `json:"users"`
}

func NewTradeResponse(trade *tradedomain.Trade) *TradeResponse {
	users := make(map[string]*TradeUserResponse, len(trade.Users))
	for id, user := range trade.Users {
		items := make([]*TradeItemResponse, len(user.Items))
		for i, item := range user.Items {
			items[i] = &TradeItemResponse{
				Type: item.Type.String(),
			}
			if item.Type == tradedomain.RoosterTradeType {
				rooster := item.Rooster()
				items[i].Data = &UserRoosterResponse{
					ID:     rooster.ID.String(),
					UserID: user.ID.String(),
					Origin: rooster.Origin,
					Type:   rooster.Type,
				}
			} else if item.Type == tradedomain.PetTradeType {
				petEntity := item.Pet()
				items[i].Data = &UserPetResponse{
					ID:     petEntity.ID.String(),
					UserID: user.ID.String(),
					Type:   int(petEntity.Type),
				}
			} else {
				itemEntity := item.Item()
				items[i].Data = &UserItemResponse{
					ID:       itemEntity.ID.String(),
					UserID:   user.ID.String(),
					Quantity: itemEntity.Quantity,
					ItemID:   itemEntity.ItemID,
					Type:     int(itemEntity.Type),
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
		Type:  "trade_update",
		ID:    trade.ID.String(),
		Users: users,
	}
}

type TradeConfirmedResponse struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func NewTradeConfirmedResponse(tradeID uuid.UUID) *TradeConfirmedResponse {
	return &TradeConfirmedResponse{
		Type: "trade_confirmed",
		ID:   tradeID.String(),
	}
}

type StartCountdownResponse struct {
	Type      string `json:"type"`
	TradeID   string `json:"trade_id"`
	Countdown int    `json:"countdown"`
}

func NewStartCountdownResponse(tradeID uuid.UUID, countdown int) *StartCountdownResponse {
	return &StartCountdownResponse{
		Type:      "start_countdown",
		TradeID:   tradeID.String(),
		Countdown: countdown,
	}
}

type TradeErrorResponse struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Error string `json:"error"`
}

func NewTradeErrorResponse(tradeID uuid.UUID, err string) *TradeErrorResponse {
	return &TradeErrorResponse{
		Type:  "trade_error",
		ID:    tradeID.String(),
		Error: err,
	}
}
