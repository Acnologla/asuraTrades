package service

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/port"
)

type UserService struct {
	userRepository    port.UserRepository
	roosterRepository port.RoosterRepository
	itemRepository    port.ItemRepository
}

func (s *UserService) Get(ctx context.Context, id domain.ID) (*domain.User, error) {
	return s.userRepository.Get(ctx, id)
}

func (s *UserService) GetUserProfile(ctx context.Context, id domain.ID) (*domain.UserProfile, error) {
	user, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.itemRepository.GetUserItems(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	roosters, err := s.roosterRepository.GetUserRoosters(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return domain.NewUserProfile(user, roosters, items), nil
}

func NewUserService(userRepository port.UserRepository, roosterRepository port.RoosterRepository, itemRepository port.ItemRepository) *UserService {
	return &UserService{
		userRepository:    userRepository,
		roosterRepository: roosterRepository,
		itemRepository:    itemRepository,
	}
}
