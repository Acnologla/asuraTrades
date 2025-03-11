package controllers

import (
	"net/http"
	"strconv"

	"github.com/acnologla/asuraTrades/internal/adapters/http/dto"
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/service"
	"github.com/gin-gonic/gin"
)

type GenerateTokenRequest struct {
	ID string `json:"id"`
}

type UserTokenController struct {
	generateTokenPassword string
	userTokenService      *service.UserTokenService
}

func (u *UserTokenController) GetUserProfile(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.String(http.StatusBadRequest, "Token not found")
		return
	}
	userProfile, err := u.userTokenService.GetUserProfile(c, token)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.NewUserProfileDTO(userProfile))
}

func (u *UserTokenController) GenerateToken(c *gin.Context) {
	password := c.Request.Header.Get("Password")
	if password != u.generateTokenPassword {
		c.String(http.StatusNotFound, "Unautorized")
		return
	}

	requestData := &GenerateTokenRequest{}
	if err := c.ShouldBindJSON(requestData); err != nil {
		c.String(http.StatusBadRequest, "ID not found")
		return
	}

	id, err := strconv.ParseUint(requestData.ID, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID format")
		return
	}

	token, err := u.userTokenService.CreateToken(c, domain.ID(id))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func NewUserTokenController(generateTokenPassword string, userTokenService *service.UserTokenService) *UserTokenController {
	return &UserTokenController{generateTokenPassword: generateTokenPassword, userTokenService: userTokenService}
}
