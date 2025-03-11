package http

import (
	"fmt"

	"github.com/acnologla/asuraTrades/internal/adapters/config"
	"github.com/acnologla/asuraTrades/internal/adapters/http/controllers"
	"github.com/gin-gonic/gin"
)

func CreateAndServe(c config.HTTPConfig, userToken *controllers.UserTokenController) error {
	r := gin.New()

	r.POST("/token", userToken.GenerateToken)

	return r.Run(fmt.Sprintf(":%s", c.Port))
}
