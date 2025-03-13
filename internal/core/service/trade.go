package service

import (
	"context"
	"errors"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/google/uuid"
)

type TradeService struct {
	cache       port.TradeCache
	userService *UserService
}

func (s *TradeService) GetTrade(ctx context.Context, id uuid.UUID) (*domain.Trade, error) {
	return s.cache.Get(id)
}

func (s *TradeService) CreateTrade(ctx context.Context, tradeID uuid.UUID, author, other domain.ID) (*domain.Trade, error) {

	if exists, err := s.cache.Get(tradeID); err == nil && exists != nil {
		return nil, errors.New("trade already exists")
	}

	trade := domain.NewTrade(tradeID, author, other)
	if err := s.cache.Update(tradeID, trade); err != nil {
		return nil, err
	}

	return trade, nil
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

func (s *TradeService) saveAndReturn(tradeID uuid.UUID, trade *domain.Trade) (*domain.Trade, error) {
	if err := s.cache.Update(tradeID, trade); err != nil {
		return nil, err
	}
	return trade, nil
}

func (s *TradeService) UpdateItem(ctx context.Context, tradeID uuid.UUID, item *dto.TradeItemDTO) (*domain.Trade, error) {
	trade, err := s.cache.Get(tradeID)
	if err != nil {
		return nil, err
	}

	if item.Remove {
		if err := trade.RemoveItem(item.User, item.ID); err != nil {
			return nil, err
		}
		return s.saveAndReturn(tradeID, trade)
	}

	tradeItem, err := s.getUserItem(ctx, item)
	if err != nil {
		return nil, err
	}

	if err := trade.AddItem(item.User, tradeItem); err != nil {
		return nil, err
	}

	return s.saveAndReturn(tradeID, trade)
}

func NewTradeService(cache port.TradeCache, userService *UserService) *TradeService {
	return &TradeService{
		cache:       cache,
		userService: userService,
	}
}
