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

type JwtTokenService struct {
	secret []byte
}

func (j *JwtTokenService) GenerateToken(id domain.ID, minutesToExpire int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": id,
			"exp": time.Now().Add(time.Minute * time.Duration(minutesToExpire)).Unix(),
		})

	return token.SignedString(j.secret)
}

func (j *JwtTokenService) ValidateToken(tokenStr string) (domain.ID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return domain.ID(claims["sub"].(float64)), nil
	}

	return 0, errors.New("token without id")
}

func NewJwtTokenService(config config.JWTConfig) port.TokenProvider {
	return &JwtTokenService{secret: []byte(config.Secret)}
}
