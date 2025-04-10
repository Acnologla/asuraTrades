package service

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/domain/trade"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository    port.UserRepository
	roosterRepository port.RoosterRepository
	itemRepository    port.ItemRepository
	petRepository     port.PetRepository
}

func (s *UserService) Get(ctx context.Context, id domain.ID) (*domain.User, error) {
	return s.userRepository.Get(ctx, id)
}

func (s *UserService) GetItem(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	return s.itemRepository.Get(ctx, id)
}

func (s *UserService) GetRooster(ctx context.Context, id uuid.UUID) (*domain.Rooster, error) {
	return s.roosterRepository.Get(ctx, id)
}

func (s *UserService) GetPet(ctx context.Context, id uuid.UUID) (*domain.Pet, error) {
	return s.petRepository.Get(ctx, id)
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

	pets, err := s.petRepository.GetUserPets(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return domain.NewUserProfile(user, trade.GetTradableEntities(roosters), trade.GetTradableEntities(items), trade.GetTradableEntities(pets)), nil
}

func NewUserService(userRepository port.UserRepository, roosterRepository port.RoosterRepository, itemRepository port.ItemRepository, petRepository port.PetRepository) *UserService {
	return &UserService{
		userRepository:    userRepository,
		roosterRepository: roosterRepository,
		itemRepository:    itemRepository,
		petRepository:     petRepository,
	}
}
