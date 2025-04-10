package http

import (
	"fmt"

	"github.com/acnologla/asuraTrades/internal/adapters/config"
	"github.com/acnologla/asuraTrades/internal/adapters/http/controllers"
	"github.com/acnologla/asuraTrades/internal/adapters/http/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func generateCorsConfig(domain string) cors.Config {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{fmt.Sprintf("https://%s", domain)}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	return config
}

func CreateAndServe(c config.HTTPConfig, userToken *controllers.UserTokenController, wsController *websocket.TradeWebsocket) error {
	r := gin.New()

	if c.Production {
		r.Use(cors.New(generateCorsConfig(c.ProductionURL)))
	}

	r.GET("/user/:token", userToken.GetUserProfile)
	r.POST("/token", userToken.GenerateToken)
	r.GET("/ws", wsController.UpgradeConnection)

	return r.Run(fmt.Sprintf(":%s", c.Port))
}
