package service

import "github.com/acnologla/asuraTrades/internal/core/port"

type TradeService struct {
	cache             port.TradeCache
	userRepository    port.UserRepository
	roosterRepository port.RoosterRepository
	itemReposotory    port.ItemRepository
}

func NewTradeService(cache port.TradeCache, userRepository port.UserRepository, roosterRepository port.RoosterRepository, itemRepository port.ItemRepository) *TradeService {
	return &TradeService{
		cache:             cache,
		userRepository:    userRepository,
		roosterRepository: roosterRepository,
		itemReposotory:    itemRepository,
	}
}
