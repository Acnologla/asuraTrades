package websocket

import (
	"context"
	"errors"
	"sync"

	"github.com/acnologla/asuraTrades/internal/adapters/http/response"
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RoomMessageType int

const (
	_ RoomMessageType = iota
	UpdateItem
	UpdateUserStatus
)

type RoomMessageData struct {
	ItemID    uuid.UUID `json:"item_id"`
	Remove    bool      `json:"remove"`
	Confirmed bool      `json:"confirmed"`
	Type      int       `json:"type"`
}

type RoomMessage struct {
	Type    RoomMessageType `json:"type"`
	TradeID uuid.UUID       `json:"trade_id"`
	User    domain.ID       `json:"user"`
	Data    RoomMessageData `json:"data"`
}

type TradeRoom struct {
	ID uuid.UUID
	sync.RWMutex
	users  map[domain.ID]*websocket.Conn
	Cancel context.CancelFunc
}

func (t *TradeRoom) AddUser(user *domain.UserTrade, connection *websocket.Conn) {
	t.users[user.AuthorID] = connection
}

func (t *TradeRoom) RemoveUser(user domain.ID) {
	delete(t.users, user)
}

func (t *TradeRoom) Broadcast(v any) {
	for user, conn := range t.users {
		if err := conn.WriteJSON(v); err != nil {
			t.RemoveUser(user)
		}
	}
}

func (t *TradeRoom) UpdateTrade(trade *domain.Trade) {
	tradeResponse := response.NewTradeResponse(trade)
	t.Broadcast(tradeResponse)
}

func NewTradeRoom(connection *websocket.Conn, tradeUser *domain.UserTrade) *TradeRoom {
	return &TradeRoom{
		ID: tradeUser.TradeID,
		users: map[domain.ID]*websocket.Conn{
			tradeUser.AuthorID: connection,
		},
	}
}

var rooms = make(map[uuid.UUID]*TradeRoom)

func GetOrCreateRoom(connection *websocket.Conn, tradeUser *domain.UserTrade) *TradeRoom {
	if room, ok := rooms[tradeUser.TradeID]; ok {
		room.Lock()
		defer room.Unlock()
		room.AddUser(tradeUser, connection)
		return room
	}

	room := NewTradeRoom(connection, tradeUser)
	rooms[tradeUser.TradeID] = room
	return room
}

func RoomMessageToTradeItemDTO(message *RoomMessage) *dto.TradeItemDTO {
	return dto.NewTradeItemDTO(message.Data.Type, message.TradeID, message.Data.ItemID, message.User, message.Data.Remove)
}

func RoomMessageToUpdateUserStatusDTO(message *RoomMessage) *dto.UpdateUserStatusDTO {
	return dto.NewUpdateUserStatusDTO(message.TradeID, message.Data.Confirmed, message.User)
}

func RemoveUserFromRoom(tradeID uuid.UUID, user domain.ID) {
	if room, ok := rooms[tradeID]; ok {
		room.Lock()
		room.RemoveUser(user)
		room.Unlock()
		if len(room.users) == 0 {
			delete(rooms, tradeID)
		}
	}
}

func validateRoomMessage(message *RoomMessage) error {
	if message.Type != UpdateItem && message.Type != UpdateUserStatus {
		return errors.New("invalid message type")
	}

	if message.TradeID == uuid.Nil {
		return errors.New("trade ID cannot be empty")
	}

	if message.Type == UpdateItem && message.Data.ItemID == uuid.Nil {
		return errors.New("item ID cannot be empty for UpdateItem message")
	}

	return nil
}
