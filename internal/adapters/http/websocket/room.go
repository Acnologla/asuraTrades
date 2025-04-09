package websocket

import (
	"context"
	"errors"
	"sync"

	"github.com/acnologla/asuraTrades/internal/adapters/http/response"
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/domain/trade"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type roomMessageType int

const (
	_ roomMessageType = iota
	UpdateItem
	UpdateUserStatus
)

type roomMessageData struct {
	ItemID    uuid.UUID `json:"item_id"`
	Remove    bool      `json:"remove"`
	Confirmed bool      `json:"confirmed"`
	Type      int       `json:"type"`
}

type roomMessage struct {
	Type    roomMessageType `json:"type"`
	TradeID uuid.UUID       `json:"trade_id"`
	User    domain.ID       `json:"user"`
	Data    roomMessageData `json:"data"`
}

type tradeRoom struct {
	id uuid.UUID
	sync.RWMutex
	users  map[domain.ID]*websocket.Conn
	cancel context.CancelFunc
}

func (t *tradeRoom) addUser(user *domain.UserTrade, connection *websocket.Conn) {
	t.users[user.AuthorID] = connection
}

func (t *tradeRoom) removeUser(user domain.ID) {
	delete(t.users, user)
}

func (t *tradeRoom) broadcast(v any) {
	for user, conn := range t.users {
		if err := conn.WriteJSON(v); err != nil {
			t.removeUser(user)
		}
	}
}

func (t *tradeRoom) updateTrade(trade *trade.Trade) {
	tradeResponse := response.NewTradeResponse(trade)
	t.broadcast(tradeResponse)
}

func newTradeRoom(connection *websocket.Conn, tradeUser *domain.UserTrade) *tradeRoom {
	return &tradeRoom{
		id: tradeUser.TradeID,
		users: map[domain.ID]*websocket.Conn{
			tradeUser.AuthorID: connection,
		},
	}
}

var rooms = make(map[uuid.UUID]*tradeRoom)

func getOrCreateRoom(connection *websocket.Conn, tradeUser *domain.UserTrade) *tradeRoom {
	if room, ok := rooms[tradeUser.TradeID]; ok {
		room.Lock()
		defer room.Unlock()
		room.addUser(tradeUser, connection)
		return room
	}

	room := newTradeRoom(connection, tradeUser)
	rooms[tradeUser.TradeID] = room
	return room
}

func roomMessageToTradeItemDTO(message *roomMessage) *dto.TradeItemDTO {
	return dto.NewTradeItemDTO(message.Data.Type, message.TradeID, message.Data.ItemID, message.User, message.Data.Remove)
}

func roomMessageToUpdateUserStatusDTO(message *roomMessage) *dto.UpdateUserStatusDTO {
	return dto.NewUpdateUserStatusDTO(message.TradeID, message.Data.Confirmed, message.User)
}

func removeUserFromRoom(tradeID uuid.UUID, user domain.ID) {
	if room, ok := rooms[tradeID]; ok {
		room.Lock()
		room.removeUser(user)
		room.Unlock()
		if len(room.users) == 0 {
			delete(rooms, tradeID)
		}
	}
}

func validateRoomMessage(message *roomMessage) error {
	if message.Type != UpdateItem && message.Type != UpdateUserStatus {
		return errors.New("invalid message type")
	}

	if message.TradeID == uuid.Nil {
		return errors.New("trade ID cannot be empty")
	}

	if message.Type == UpdateItem {
		if message.Data.ItemID == uuid.Nil {
			return errors.New("item ID cannot be empty for UpdateItem message")
		}
	}

	return nil
}
