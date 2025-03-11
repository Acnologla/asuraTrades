package service

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/port"
)

type UserTokenService struct {
	tokenProvider     port.TokenProvider
	userRepository    port.UserRepository
	itemRepository    port.ItemRepository
	roosterRepository port.RoosterRepository
}

func (s *UserTokenService) CreateToken(ctx context.Context, id domain.ID) (string, error) {
	_, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return "", err
	}

	return s.tokenProvider.GenerateToken(id, 20)
}

func (s *UserTokenService) GetUserProfile(ctx context.Context, token string) (*domain.UserProfile, error) {
	userID, err := s.tokenProvider.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	items, err := s.itemRepository.GetUserItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	roosters, err := s.roosterRepository.GetUserRoosters(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserProfile{
		User:     user,
		Items:    items,
		Roosters: roosters,
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
