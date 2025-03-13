package service

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository    port.UserRepository
	roosterRepository port.RoosterRepository
	itemRepository    port.ItemRepository
}

func (s *UserService) Get(ctx context.Context, id domain.ID) (*domain.User, error) {
	return s.userRepository.Get(ctx, id)
}

func (s *UserService) Lock(ctx context.Context, id domain.ID) (func(error) error, error) {
	return s.userRepository.LockUpdate(ctx, id)
}

func (s *UserService) GetItem(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	return s.itemRepository.Get(ctx, id)
}

func (s *UserService) GetRooster(ctx context.Context, id uuid.UUID) (*domain.Rooster, error) {
	return s.roosterRepository.Get(ctx, id)
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

	return domain.NewUserProfile(user, domain.GetTradableRoosters(roosters), domain.GetTradableItems(items)), nil
}

func NewUserService(userRepository port.UserRepository, roosterRepository port.RoosterRepository, itemRepository port.ItemRepository) *UserService {
	return &UserService{
		userRepository:    userRepository,
		roosterRepository: roosterRepository,
		itemRepository:    itemRepository,
	}
}
