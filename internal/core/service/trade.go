package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/google/uuid"
)

type UpdateUserStatusWrapper struct {
	Trade *domain.Trade
	Done  bool
}

type RoosterTransferRequest struct {
	Rooster        *domain.Rooster
	NewOwnerID     domain.ID
	CurrentOwnerID domain.ID
}

type ItemTransferRequest struct {
	Item       *domain.Item
	NewOwnerID domain.ID
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

func (s *TradeService) transferItem(ctx context.Context, request *ItemTransferRequest, repo port.ItemRepository) error {
	i, err := repo.Get(ctx, request.Item.ID)
	if err != nil {
		return err
	}

	if i.Quantity < request.Item.Quantity {
		return errors.New("not enough items")
	}

	if err := repo.Remove(ctx, request.Item.ID, request.Item.Quantity); err != nil {
		return err
	}

	newItem := domain.NewItem(request.NewOwnerID, request.Item.ItemID, request.Item.Type)
	return repo.Add(ctx, newItem, request.Item.Quantity)
}

func (s *TradeService) transferRooster(ctx context.Context, request *RoosterTransferRequest, repo port.RoosterRepository) error {
	if _, err := repo.Get(ctx, request.Rooster.ID); err != nil {
		return err
	}

	if err := repo.Delete(ctx, request.Rooster.ID); err != nil {
		return err
	}

	origin := fmt.Sprintf("Trade with %s", request.CurrentOwnerID)
	newRooster := domain.NewRooster(request.NewOwnerID, request.Rooster.Type, origin)
	return repo.Create(ctx, newRooster)
}

func (s *TradeService) swapUserItems(ctx context.Context, user *domain.TradeUser, otherID domain.ID, adapters port.UserTradeTxAdapters) error {
	for _, item := range user.Items {
		if item.Type == domain.ItemTradeType {
			if err := s.transferItem(ctx, &ItemTransferRequest{Item: item.Item, NewOwnerID: otherID}, adapters.ItemRepository); err != nil {
				return err
			}
		} else {
			if err := s.transferRooster(ctx, &RoosterTransferRequest{Rooster: item.Rooster, CurrentOwnerID: user.ID, NewOwnerID: otherID}, adapters.RoosterRepository); err != nil {
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

const COUNTDOWN_SECONDS = 5

func (s *TradeService) ConfirmTrade(ctx context.Context, tradeID uuid.UUID, callback func(bool, error)) (int, error) {
	trade, err := s.GetTrade(ctx, tradeID)
	if err != nil {
		return 0, err
	}

	if !trade.Done() {
		return 0, errors.New("trade not done")
	}

	go func() {
		ticker := time.NewTicker(COUNTDOWN_SECONDS * time.Second)
		defer ticker.Stop()
		select {
		case <-ticker.C:
			err := s.FinishTrade(ctx, tradeID)
			callback(err == nil, err)
		case <-ctx.Done():
			callback(false, nil)
			return
		}
	}()

	return COUNTDOWN_SECONDS, nil
}

func (s *TradeService) UpdateUserStatus(ctx context.Context, dto *dto.UpdateUserStatusDTO) (*UpdateUserStatusWrapper, error) {
	trade, err := s.GetTrade(ctx, dto.ID)
	if err != nil {
		return nil, err
	}

	if err := trade.UpdateUserStatus(dto.User, dto.Confirmed); err != nil {
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

func (s *TradeService) UpdateItem(ctx context.Context, dto *dto.TradeItemDTO) (*domain.Trade, error) {
	trade, err := s.GetTrade(ctx, dto.ID)
	if err != nil {
		return nil, err
	}

	if trade.Users[dto.User].Confirmed {
		return nil, errors.New("user already confirmed")
	}

	if dto.Remove {
		if err := trade.RemoveItem(dto.User, dto.ItemID, dto.Type); err != nil {
			return nil, err
		}
		return s.saveAndReturn(dto.ID, trade)
	}

	tradeItem, err := s.getUserItem(ctx, dto.ItemID, dto.Type)
	if err != nil {
		return nil, err
	}

	if err := trade.AddItem(dto.User, tradeItem); err != nil {
		return nil, err
	}

	return s.saveAndReturn(dto.ID, trade)
}

func NewTradeService(cache port.TradeCache, userService *UserService, userTx port.TradeTxProvider) *TradeService {
	return &TradeService{
		cache:       cache,
		userService: userService,
		userTx:      userTx,
	}
}
