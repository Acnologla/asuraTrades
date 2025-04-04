package websocket

import (
	"context"
	"net/http"
	"net/url"
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

type TradeWebsocket struct {
	tokenService     *service.UserTokenService
	tradeService     *service.TradeService
	productionDomain string
	production       bool
}

func (t *TradeWebsocket) getUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if !t.production {
				return true
			}
			origin := r.Header.Get("Origin")
			originURL, err := url.Parse(origin)
			if err != nil {
				return false
			}

			return originURL.Host == t.productionDomain
		},
	}
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

func (t *TradeWebsocket) confirmTrade(ctx context.Context, room *tradeRoom) {
	context, fn := context.WithCancel(ctx)
	room.cancel = fn
	countdownTime, err := t.tradeService.ConfirmTrade(context, room.id, func(b bool, err error) {
		if b {
			room.broadcast(response.NewTradeConfirmedResponse(room.id))
			for _, conn := range room.users {
				_ = conn.Close()
			}
		}
		if err != nil {
			room.broadcast(response.NewTradeErrorResponse(room.id, err.Error()))
		}
	})

	if err != nil {
		room.cancel = nil
		return
	}

	room.broadcast(response.NewStartCountdownResponse(room.id, countdownTime))
}

func (t *TradeWebsocket) processMessage(ctx context.Context, room *tradeRoom, message *roomMessage) {
	room.Lock()
	defer room.Unlock()

	err := validateRoomMessage(message)
	if err != nil {
		return
	}

	if message.Type == UpdateItem {
		UpdateItemDTO := roomMessageToTradeItemDTO(message)
		if trade, err := t.tradeService.UpdateItem(ctx, UpdateItemDTO); err == nil {
			room.updateTrade(trade)
		}
	}
	if message.Type == UpdateUserStatus {
		UpdateUserStatusDTO := roomMessageToUpdateUserStatusDTO(message)
		if r, err := t.tradeService.UpdateUserStatus(ctx, UpdateUserStatusDTO); err == nil {
			room.updateTrade(r.Trade)

			if room.cancel != nil {
				room.cancel()
				room.cancel = nil
			}

			if r.Done {
				t.confirmTrade(ctx, room)
			}
		}
	}

}

func (t *TradeWebsocket) sendPongs(conn *websocket.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(PONG_INTERVAL))

	conn.SetPingHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(PONG_INTERVAL))
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

func (t *TradeWebsocket) initializeUser(ctx context.Context, conn *websocket.Conn, room *tradeRoom, user *domain.UserTrade) {
	defer removeUserFromRoom(user.TradeID, user.AuthorID)
	t.sendPongs(conn)
	limiter := rate.NewLimiter(1, 3)
	for {
		if !limiter.Allow() {
			continue
		}

		message := &roomMessage{}
		if err := conn.ReadJSON(message); err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				_ = conn.WriteControl(websocket.CloseMessage,
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
	conn, err := t.getUpgrader().Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error occurred while upgrading connection")
		return
	}
	defer conn.Close()
	room := getOrCreateRoom(conn, tokenInfo)
	trade, err := t.tradeService.GetTrade(c.Request.Context(), room.id)
	if err != nil {
		trade, err = t.tradeService.CreateTrade(c.Request.Context(), tokenInfo.TradeID, tokenInfo.AuthorID, tokenInfo.OtherID)
		if err != nil {
			return
		}
	}
	_ = conn.WriteJSON(response.NewTradeResponse(trade))

	t.initializeUser(c.Request.Context(), conn, room, tokenInfo)
}

func NewTradeWebsocket(tokenService *service.UserTokenService, tradeService *service.TradeService, prooduction bool, productionURl string) *TradeWebsocket {
	return &TradeWebsocket{
		tokenService: tokenService,
		tradeService: tradeService,
	}
}
