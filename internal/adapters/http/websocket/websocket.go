package websocket

import (
	"context"
	"net/http"
	"time"

	"github.com/acnologla/asuraTrades/internal/adapters/http/response"
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

const PING_INTERVAL = 30 * time.Second
const PONG_INTERVAL = 60 * time.Second

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type TradeWebsocket struct {
	tokenService *service.UserTokenService
	tradeService *service.TradeService
}

func (t *TradeWebsocket) authAndDecodeToken(c *gin.Context) *domain.UserTrade {
	token := c.Query("token")
	if token == "" {
		c.String(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return nil
	}

	trade, err := t.tokenService.DecodeToken(token)
	if err != nil {
		c.String(http.StatusForbidden, "Forbidden")
		c.Abort()
		return nil
	}

	return trade
}

func (t *TradeWebsocket) confirmTrade(ctx context.Context, room *TradeRoom) {
	context, fn := context.WithCancel(ctx)
	room.Cancel = fn
	countdownTime, err := t.tradeService.ConfirmTrade(context, room.ID, func(b bool, err error) {
		if b {
			room.Broadcast(response.NewTradeConfirmedResponse(room.ID))
			for _, conn := range room.users {
				conn.Close()
			}
		}
	})

	if err != nil {
		room.Cancel = nil
		return
	}

	room.Broadcast(response.NewStartCountdownResponse(room.ID, countdownTime))
}

func (t *TradeWebsocket) processMessage(ctx context.Context, room *TradeRoom, message *RoomMessage) {
	room.Lock()
	defer room.Unlock()

	err := validateRoomMessage(message)
	if err != nil {
		return
	}

	if message.Type == UpdateItem {
		UpdateItemDTO := RoomMessageToTradeItemDTO(message)
		if trade, err := t.tradeService.UpdateItem(ctx, UpdateItemDTO); err == nil {
			room.UpdateTrade(trade)
		}
	}
	if message.Type == UpdateUserStatus {
		UpdateUserStatusDTO := RoomMessageToUpdateUserStatusDTO(message)
		if r, err := t.tradeService.UpdateUserStatus(ctx, UpdateUserStatusDTO); err == nil {
			room.UpdateTrade(r.Trade)

			if room.Cancel != nil {
				room.Cancel()
				room.Cancel = nil
			}

			if r.Done {
				t.confirmTrade(ctx, room)
			}
		}
	}

}

func (t *TradeWebsocket) sendPongs(conn *websocket.Conn) {
	conn.SetReadDeadline(time.Now().Add(PONG_INTERVAL))

	conn.SetPingHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(PONG_INTERVAL))
		return conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(10*time.Second))
	})

	go func() {
		ticker := time.NewTicker(PING_INTERVAL)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				return
			}
		}
	}()
}

func (t *TradeWebsocket) initializeUser(ctx context.Context, conn *websocket.Conn, room *TradeRoom, user *domain.UserTrade) {
	defer RemoveUserFromRoom(user.TradeID, user.AuthorID)
	t.sendPongs(conn)
	limiter := rate.NewLimiter(1, 3)
	for {
		if !limiter.Allow() {
			continue
		}

		message := &RoomMessage{}
		if err := conn.ReadJSON(message); err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				conn.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Goodbye"),
					time.Now().Add(time.Second))
			}
			return
		}
		message.User = user.AuthorID
		t.processMessage(ctx, room, message)
	}
}

func (t *TradeWebsocket) UpgradeConnection(c *gin.Context) {
	tokenInfo := t.authAndDecodeToken(c)
	if tokenInfo == nil {
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error occurred while upgrading connection")
		return
	}
	defer conn.Close()
	room := GetOrCreateRoom(conn, tokenInfo)
	trade, err := t.tradeService.GetTrade(c.Request.Context(), room.ID)
	if err != nil {
		trade, err = t.tradeService.CreateTrade(c.Request.Context(), tokenInfo.TradeID, tokenInfo.AuthorID, tokenInfo.OtherID)
		if err != nil {
			return
		}
	}
	conn.WriteJSON(response.NewTradeResponse(trade))

	t.initializeUser(c.Request.Context(), conn, room, tokenInfo)
}

func NewTradeWebsocket(tokenService *service.UserTokenService, tradeService *service.TradeService) *TradeWebsocket {
	return &TradeWebsocket{
		tokenService: tokenService,
		tradeService: tradeService,
	}
}
