package service

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/acnologla/asuraTrades/internal/core/port"
)

type UserTokenService struct {
	tokenProvider     port.TokenProvider
	userRepository    port.UserRepository
	itemRepository    port.ItemRepository
	roosterRepository port.RoosterRepository
}

type GetTradeTokenResponseWrapper struct {
	UserTrade   *domain.UserTrade
	UserProfile *domain.UserProfile
}

func (s *UserTokenService) CreateToken(ctx context.Context, userTradeDto *dto.GenerateUserTokenDTO) (string, error) {
	userTrade, err := domain.NewUserTrade(userTradeDto.AuthorID, userTradeDto.OtherID, userTradeDto.TradeID)
	if err != nil {
		return "", err
	}
	_, err = s.userRepository.Get(ctx, userTrade.AuthorID)
	if err != nil {
		return "", err
	}

	_, err = s.userRepository.Get(ctx, userTrade.OtherID)
	if err != nil {
		return "", err
	}

	return s.tokenProvider.GenerateToken(userTrade, 20)
}

func (s *UserTokenService) GetTradeTokenResponse(ctx context.Context, token string) (*GetTradeTokenResponseWrapper, error) {
	userTrade, err := s.tokenProvider.ValidateToken(token)
	authorID := userTrade.AuthorID
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.Get(ctx, authorID)
	if err != nil {
		return nil, err
	}

	items, err := s.itemRepository.GetUserItems(ctx, authorID)
	if err != nil {
		return nil, err
	}

	roosters, err := s.roosterRepository.GetUserRoosters(ctx, authorID)
	if err != nil {
		return nil, err
	}

	return &GetTradeTokenResponseWrapper{
		UserTrade: userTrade,
		UserProfile: &domain.UserProfile{
			User:     user,
			Items:    items,
			Roosters: roosters,
		},
	}, nil
}

func NewUserTokenService(tp port.TokenProvider, ur port.UserRepository, ir port.ItemRepository, rr port.RoosterRepository) *UserTokenService {
	return &UserTokenService{
		tokenProvider:     tp,
		userRepository:    ur,
		itemRepository:    ir,
		roosterRepository: rr,
	}
}
