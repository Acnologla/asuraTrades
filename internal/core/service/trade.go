package service

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/google/uuid"
)

type TradeService struct {
	cache       port.TradeCache
	userService *UserService
}

func (s *TradeService) getUserItem(ctx context.Context, dto *dto.TradeItemDTO) (*domain.TradeItem, error) {
	if dto.Type == domain.ItemTradeType {
		i, err := s.userService.GetItem(ctx, dto.ID)
		if err != nil {
			return nil, err
		}
		return domain.NewTradeItemItem(i), nil
	}
	r, err := s.userService.GetRooster(ctx, dto.ID)
	if err != nil {
		return nil, err
	}
	return domain.NewTradeItemRooster(r), nil
}

func (s *TradeService) AddItem(ctx context.Context, tradeID uuid.UUID, item *dto.TradeItemDTO) (*domain.Trade, error) {
	trade, err := s.cache.Get(tradeID)
	if err != nil {
		return nil, err
	}

	tradeItem, err := s.getUserItem(ctx, item)
	if err != nil {
		return nil, err
	}

	if err := trade.AddItem(item.User, tradeItem); err != nil {
		return nil, err
	}

	if err := s.cache.Update(tradeID, trade); err != nil {
		return nil, err
	}

	return trade, nil
}

func NewTradeService(cache port.TradeCache, userService *UserService) *TradeService {
	return &TradeService{
		cache:       cache,
		userService: userService,
	}
}
