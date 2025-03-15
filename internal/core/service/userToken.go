package service

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/acnologla/asuraTrades/internal/core/port"
)

type UserTokenService struct {
	tokenProvider port.TokenProvider
	userService   *UserService
}

type GetTradeTokenResponseWrapper struct {
	UserTrade   *domain.UserTrade
	UserProfile *domain.UserProfile
}

const TOKEN_EXPIRATION_TIME = 25

func (s *UserTokenService) CreateToken(ctx context.Context, userTradeDto *dto.GenerateUserTokenDTO) (string, error) {
	userTrade, err := domain.NewUserTrade(userTradeDto.AuthorID, userTradeDto.OtherID, userTradeDto.TradeID)
	if err != nil {
		return "", err
	}
	_, err = s.userService.Get(ctx, userTrade.AuthorID)
	if err != nil {
		return "", err
	}

	_, err = s.userService.Get(ctx, userTrade.OtherID)
	if err != nil {
		return "", err
	}

	return s.tokenProvider.GenerateToken(userTrade, TOKEN_EXPIRATION_TIME)
}

func (s *UserTokenService) DecodeToken(token string) (*domain.UserTrade, error) {
	return s.tokenProvider.ValidateToken(token)
}

func (s *UserTokenService) GetTradeTokenResponse(ctx context.Context, token string) (*GetTradeTokenResponseWrapper, error) {
	userTrade, err := s.tokenProvider.ValidateToken(token)
	authorID := userTrade.AuthorID
	if err != nil {
		return nil, err
	}

	profile, err := s.userService.GetUserProfile(ctx, authorID)
	if err != nil {
		return nil, err
	}

	return &GetTradeTokenResponseWrapper{
		UserTrade:   userTrade,
		UserProfile: profile,
	}, nil
}

func NewUserTokenService(tp port.TokenProvider, us *UserService) *UserTokenService {
	return &UserTokenService{
		tokenProvider: tp,
		userService:   us,
	}
}
