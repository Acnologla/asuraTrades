package websocket

import (
	"sync"

	"github.com/acnologla/asuraTrades/internal/core/domain"
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
	users map[domain.ID]*websocket.Conn
}

func (t *TradeRoom) AddUser(user *domain.UserTrade, connection *websocket.Conn) {
	t.Lock()
	defer t.Unlock()
	t.users[user.AuthorID] = connection
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
		room.AddUser(tradeUser, connection)
		return room
	}

	room := NewTradeRoom(connection, tradeUser)
	rooms[tradeUser.TradeID] = room
	return room
}
