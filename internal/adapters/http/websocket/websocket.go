package websocket

import (
	"net/http"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type TradeWebsocket struct {
	tokenService service.UserTokenService
	tradeService service.TradeService
}

func (t *TradeWebsocket) authAndDecodeToken(c *gin.Context) *domain.UserTrade {
	token := c.Request.Header.Get("Authorization")
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

func (t *TradeWebsocket) UpgradeConnection(c *gin.Context) {
	token := t.authAndDecodeToken(c)
	if token == nil {
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error occurred while upgrading connection")
		return
	}
	defer conn.Close()
	room := GetOrCreateRoom(conn, token)
	trade, err := t.tradeService.GetTrade(c.Request.Context(), room.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error occurred while getting trade")
		return
	}

	conn.WriteJSON(trade)

}

func NewTradeWebsocket(tokenService service.UserTokenService, tradeService service.TradeService) *TradeWebsocket {
	return &TradeWebsocket{
		tokenService: tokenService,
		tradeService: tradeService,
	}
}
