package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/acnologla/asuraTrades/internal/adapters/config"
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenProvider struct {
	secret []byte
}

func (j *JwtTokenProvider) GenerateToken(trade *domain.UserTrade, minutesToExpire int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":     trade.AuthorID.String(),
			"otherID": trade.OtherID.String(),
			"tradeID": trade.TradeID.String(),
			"exp":     time.Now().Add(time.Minute * time.Duration(minutesToExpire)).Unix(),
		})

	return token.SignedString(j.secret)
}

func (j *JwtTokenProvider) ValidateToken(tokenStr string) (*domain.UserTrade, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		sub, otherID, tradeID := claims["sub"].(string), claims["otherID"].(string), claims["tradeID"].(string)
		return domain.NewUserTrade(sub, otherID, tradeID)
	}

	return nil, errors.New("token without id")
}

func NewJwtTokenService(config config.JWTConfig) port.TokenProvider {
	return &JwtTokenProvider{secret: []byte(config.Secret)}
}
