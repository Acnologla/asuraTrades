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

	if exists, err := s.GetTrade(ctx, tradeID); err == nil && exists != nil {
		return nil, errors.New("trade already exists")
	}

	trade := domain.NewTrade(tradeID, author, other)
	if err := s.cache.Update(tradeID, trade); err != nil {
		return nil, err
	}

	return trade, nil
}

func (s *TradeService) transferItem(ctx context.Context, item *domain.Item, newOwnerID domain.ID, repo port.ItemRepository) error {
	if _, err := repo.Get(ctx, item.ID); err != nil {
		return err
	}

	if err := repo.Remove(ctx, item.ID); err != nil {
		return err
	}

	newItem := domain.NewItem(newOwnerID, item.ItemID, item.Type)
	return repo.Add(ctx, newItem)
}

func (s *TradeService) transferRooster(ctx context.Context, rooster *domain.Rooster, newOwnerID domain.ID, currentOwnerID domain.ID, repo port.RoosterRepository) error {
	if _, err := repo.Get(ctx, rooster.ID); err != nil {
		return err
	}

	if err := repo.Delete(ctx, rooster.ID); err != nil {
		return err
	}

	origin := fmt.Sprintf("Trade with %s", currentOwnerID)
	newRooster := domain.NewRooster(newOwnerID, rooster.Type, origin)
	return repo.Create(ctx, newRooster)
}

func (s *TradeService) swapUserItems(ctx context.Context, user *domain.TradeUser, otherID domain.ID, adapters port.UserTradeTxAdapters) error {
	for _, item := range user.Items {
		if item.Type == domain.ItemTradeType {
			if err := s.transferItem(ctx, item.Item, otherID, adapters.ItemRepository); err != nil {
				return err
			}
		} else {
			if err := s.transferRooster(ctx, item.Rooster, otherID, user.ID, adapters.RoosterRepository); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *TradeService) FinishTrade(ctx context.Context, tradeID uuid.UUID) error {
	trade, err := s.GetTrade(ctx, tradeID)
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
	trade, err := s.GetTrade(ctx, dto.ID)
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
	trade, err := s.GetTrade(ctx, tradeID)
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
