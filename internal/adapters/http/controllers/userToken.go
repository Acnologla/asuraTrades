package controllers

import (
	"github.com/gin-gonic/gin"
)

type UserTokenController struct {
	generateTokenPassword string
}

func (u *UserTokenController) GenerateToken(c *gin.Context) {
	password := c.Request.Header.Get("password")
	if password != u.generateTokenPassword {
		c.String(404, "Unautorized")
		return
	}

}

func NewUserTokenController(generateTokenPassword string) *UserTokenController {
	return &UserTokenController{generateTokenPassword: generateTokenPassword}
}
