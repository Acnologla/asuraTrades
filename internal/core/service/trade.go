package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/google/uuid"
)

type UpdateUserStatusWrapper struct {
	Trade *domain.Trade
	Done  bool
}

type TradeService struct {
	cache       port.TradeCache
	userService *UserService
	userTx      port.TradeTxProvider
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

func (s *TradeService) swapUserItems(ctx context.Context, user *domain.TradeUser, otherID domain.ID, adapters port.UserTradeTxAdapters) error {
	for _, item := range user.Items {
		if item.Type == domain.ItemTradeType {
			if _, err := adapters.ItemRepository.Get(ctx, item.Item.ID); err != nil {
				return err
			}

			if err := adapters.ItemRepository.Remove(ctx, item.Item.ID); err != nil {
				return err
			}

			if err := adapters.ItemRepository.Add(ctx, domain.NewItem(otherID, item.Item.ItemID, item.Item.Type)); err != nil {
				return err
			}
		} else {
			if _, err := adapters.RoosterRepository.Get(ctx, item.Rooster.ID); err != nil {
				return err
			}

			if err := adapters.RoosterRepository.Delete(ctx, item.Rooster.ID); err != nil {
				return err
			}

			origin := fmt.Sprintf("Trade with %s", user.ID)
			if err := adapters.RoosterRepository.Create(ctx, domain.NewRooster(otherID, item.Rooster.Type, origin)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *TradeService) FinishTrade(ctx context.Context, tradeID uuid.UUID) error {
	trade, err := s.cache.Get(tradeID)
	if err != nil {
		return err
	}

	if !trade.Done() {
		return errors.New("trade not done")
	}

	err = s.userTx.Transact(ctx, func(adapters port.UserTradeTxAdapters, lock func(domain.ID) error) error {
		users := []*domain.TradeUser{}
		for _, user := range trade.Users {
			if err := lock(user.ID); err != nil {
				return err
			}
			users = append(users, user)
		}

		firstUser, secondUser := users[0], users[1]
		if err := s.swapUserItems(ctx, firstUser, secondUser.ID, adapters); err != nil {
			return err
		}

		if err := s.swapUserItems(ctx, secondUser, firstUser.ID, adapters); err != nil {
			return err
		}

		return nil
	})

	if err == nil {
		s.cache.Delete(tradeID)
	}

	return err

}

func (s *TradeService) UpdateUserStatus(ctx context.Context, dto *dto.UpdateUserStatusDTO) (*UpdateUserStatusWrapper, error) {
	trade, err := s.cache.Get(dto.ID)
	if err != nil {
		return nil, err
	}

	if err := trade.UpdateUserStatus(dto.UserID, dto.Confirmed); err != nil {
		return nil, err
	}

	if _, err := s.saveAndReturn(dto.ID, trade); err != nil {
		return nil, err
	}

	return &UpdateUserStatusWrapper{
		Trade: trade,
		Done:  trade.Done(),
	}, nil
}

func (s *TradeService) getUserItem(ctx context.Context, id uuid.UUID, t domain.TradeItemType) (*domain.TradeItem, error) {
	if t == domain.ItemTradeType {
		i, err := s.userService.GetItem(ctx, id)
		if err != nil {
			return nil, err
		}
		return domain.NewTradeItemItem(i), nil
	}
	r, err := s.userService.GetRooster(ctx, id)
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

	tradeItem, err := s.getUserItem(ctx, item.ID, item.Type)
	if err != nil {
		return nil, err
	}

	if err := trade.AddItem(item.User, tradeItem); err != nil {
		return nil, err
	}

	return s.saveAndReturn(tradeID, trade)
}

func NewTradeService(cache port.TradeCache, userService *UserService, userTx port.TradeTxProvider) *TradeService {
	return &TradeService{
		cache:       cache,
		userService: userService,
		userTx:      userTx,
	}
}
