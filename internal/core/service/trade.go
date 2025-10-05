package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	tradedomain "github.com/acnologla/asuraTrades/internal/core/domain/trade"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/google/uuid"
)

type UpdateUserStatusWrapper struct {
	Trade *tradedomain.Trade
	Done  bool
}

type TransferRequest[T tradedomain.Tradeable] struct {
	Object         T
	NewOwnerID     domain.ID
	CurrentOwnerID domain.ID
}

func newTransferRequest[T tradedomain.Tradeable](object T, newOwnerID, currentOwnerID domain.ID) *TransferRequest[T] {
	return &TransferRequest[T]{
		Object:         object,
		NewOwnerID:     newOwnerID,
		CurrentOwnerID: currentOwnerID,
	}
}

type TradeService struct {
	cache       port.TradeCache
	userService *UserService
	userTx      port.TradeTxProvider
}

func (s *TradeService) GetTrade(ctx context.Context, id uuid.UUID) (*tradedomain.Trade, error) {
	return s.cache.Get(id)
}

func (s *TradeService) CreateTrade(ctx context.Context, tradeID uuid.UUID, author, other domain.ID) (*tradedomain.Trade, error) {

	if exists, err := s.GetTrade(ctx, tradeID); err == nil && exists != nil {
		return nil, errors.New("trade already exists")
	}

	trade := tradedomain.NewTrade(tradeID, author, other)
	if err := s.cache.Update(tradeID, trade); err != nil {
		return nil, err
	}

	return trade, nil
}

func (s *TradeService) transferItem(ctx context.Context, request *TransferRequest[*domain.Item], repo port.ItemRepository) error {
	i, err := repo.Get(ctx, request.Object.ID)
	if err != nil {
		return err
	}

	if i.Quantity < request.Object.Quantity {
		return errors.New("not enough items")
	}

	if err := repo.Remove(ctx, request.Object.ID, request.Object.Quantity); err != nil {
		return err
	}

	newItem := domain.NewItem(request.NewOwnerID, request.Object.ItemID, request.Object.Type)
	return repo.Add(ctx, newItem, request.Object.Quantity)
}

func (s *TradeService) transferPet(ctx context.Context, request *TransferRequest[*domain.Pet], repo port.PetRepository) error {
	if _, err := repo.Get(ctx, request.Object.ID); err != nil {
		return err
	}

	if err := repo.Delete(ctx, request.Object.ID); err != nil {
		return err
	}

	newPet := domain.NewPet(request.NewOwnerID, request.Object.Type, request.Object.Level)
	return repo.Create(ctx, newPet)
}

const MAX_ROOSTERS_QUANTITY = 24

func (s *TradeService) checkMaxRoosters(ctx context.Context, id domain.ID, repo port.RoosterRepository) error {
	roosterQuantity, err := repo.GetUserRoosterQuantity(ctx, id)
	if err != nil {
		return err
	}
	if roosterQuantity > MAX_ROOSTERS_QUANTITY {
		return errors.New("too many roosters")
	}
	return nil
}

func (s *TradeService) transferRooster(ctx context.Context, request *TransferRequest[*domain.Rooster], repo port.RoosterRepository) error {
	if _, err := repo.Get(ctx, request.Object.ID); err != nil {
		return err
	}

	if err := repo.Delete(ctx, request.Object.ID); err != nil {
		return err
	}

	origin := fmt.Sprintf("Trade with %s", request.CurrentOwnerID)
	newRooster := domain.NewRooster(request.NewOwnerID, request.Object.Type, origin, request.Object.Special)
	if err := repo.Create(ctx, newRooster); err != nil {
		return err
	}

	return s.checkMaxRoosters(ctx, request.NewOwnerID, repo)
}

func (s *TradeService) swapUserItems(ctx context.Context, user *tradedomain.TradeUser, otherID domain.ID, adapters port.UserTradeTxAdapters) error {
	for _, item := range user.Items {
		var err error
		switch item.Type {
		case tradedomain.ItemTradeType:
			err = s.transferItem(ctx, newTransferRequest(item.Item(), otherID, user.ID), adapters.ItemRepository)
		case tradedomain.RoosterTradeType:
			err = s.transferRooster(ctx, newTransferRequest(item.Rooster(), otherID, user.ID), adapters.RoosterRepository)
		case tradedomain.PetTradeType:
			err = s.transferPet(ctx, newTransferRequest(item.Pet(), otherID, user.ID), adapters.PetRepository)
		}

		if err != nil {
			return err
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
		users := []*tradedomain.TradeUser{}
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
		_ = s.cache.Delete(tradeID)
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

func (s *TradeService) getUserItem(ctx context.Context, id uuid.UUID, t tradedomain.TradeItemType) (*tradedomain.TradeItem, error) {
	switch t {
	case tradedomain.ItemTradeType:
		item, err := s.userService.GetItem(ctx, id)
		if err != nil {
			return nil, err
		}
		return tradedomain.NewTradeItemItem(item), nil
	case tradedomain.PetTradeType:
		pet, err := s.userService.GetPet(ctx, id)
		if err != nil {
			return nil, err
		}
		return tradedomain.NewTradeItemPet(pet), nil
	case tradedomain.RoosterTradeType:

		rooster, err := s.userService.GetRooster(ctx, id)
		if err != nil {
			return nil, err
		}
		return tradedomain.NewTradeItemRooster(rooster), nil
	}

	return nil, errors.New("invalid trade item type")
}

func (s *TradeService) saveAndReturn(tradeID uuid.UUID, trade *tradedomain.Trade) (*tradedomain.Trade, error) {
	if err := s.cache.Update(tradeID, trade); err != nil {
		return nil, err
	}
	return trade, nil
}

func (s *TradeService) removeItem(dto *dto.TradeItemDTO, trade *tradedomain.Trade) (*tradedomain.Trade, error) {
	if err := trade.RemoveItem(dto.User, dto.ItemID, dto.Type); err != nil {
		return nil, err
	}
	return s.saveAndReturn(dto.ID, trade)
}

func (s *TradeService) addItem(ctx context.Context, dto *dto.TradeItemDTO, trade *tradedomain.Trade) (*tradedomain.Trade, error) {
	tradeItem, err := s.getUserItem(ctx, dto.ItemID, dto.Type)
	if err != nil {
		return nil, err
	}

	if err := trade.AddItem(dto.User, tradeItem); err != nil {
		return nil, err
	}

	return s.saveAndReturn(dto.ID, trade)
}

func (s *TradeService) UpdateItem(ctx context.Context, dto *dto.TradeItemDTO) (*tradedomain.Trade, error) {
	trade, err := s.GetTrade(ctx, dto.ID)
	if err != nil {
		return nil, err
	}

	if trade.Users[dto.User].Confirmed {
		return nil, errors.New("user already confirmed")
	}

	if dto.Remove {
		return s.removeItem(dto, trade)
	}

	return s.addItem(ctx, dto, trade)
}

func NewTradeService(cache port.TradeCache, userService *UserService, userTx port.TradeTxProvider) *TradeService {
	return &TradeService{
		cache:       cache,
		userService: userService,
		userTx:      userTx,
	}
}
