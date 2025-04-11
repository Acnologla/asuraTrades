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

type RoomManager struct {
	sync.Mutex
	rooms map[uuid.UUID]*tradeRoom
}

func (rm *RoomManager) getOrCreateRoom(connection *websocket.Conn, tradeUser *domain.UserTrade) *tradeRoom {
	rm.Lock()
	defer rm.Unlock()

	if room, ok := rm.rooms[tradeUser.TradeID]; ok {
		room.addUser(tradeUser, connection)
		return room
	}

	room := newTradeRoom(connection, tradeUser)
	rm.rooms[tradeUser.TradeID] = room
	return room
}

func (rm *RoomManager) removeUser(tradeID uuid.UUID, userID domain.ID) {
	rm.Lock()
	defer rm.Unlock()

	if room, ok := rm.rooms[tradeID]; ok {
		room.removeUser(userID)
		if len(room.users) == 0 {
			delete(rm.rooms, tradeID)
		}
	}
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[uuid.UUID]*tradeRoom),
	}
}

func roomMessageToTradeItemDTO(message *roomMessage) *dto.TradeItemDTO {
	return dto.NewTradeItemDTO(message.Data.Type, message.TradeID, message.Data.ItemID, message.User, message.Data.Remove)
}

func roomMessageToUpdateUserStatusDTO(message *roomMessage) *dto.UpdateUserStatusDTO {
	return dto.NewUpdateUserStatusDTO(message.TradeID, message.Data.Confirmed, message.User)
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
