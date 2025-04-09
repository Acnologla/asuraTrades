package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/acnologla/asuraTrades/internal/core/port/mock"
	"github.com/acnologla/asuraTrades/internal/core/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type TradeServiceTestSuite struct {
	ctrl            *gomock.Controller
	ctx             context.Context
	mockUserRepo    *mock.MockUserRepository
	mockItemRepo    *mock.MockItemRepository
	mockRoosterRepo *mock.MockRoosterRepository
	mockTxProvider  *mock.MockTradeTxProvider
	mockCache       *mock.MockTradeCache
	userService     *service.UserService
	tradeService    *service.TradeService
}

func SetupTradeServiceTest(t *testing.T) *TradeServiceTestSuite {
	ctrl := gomock.NewController(t)
	ctx := t.Context()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	mockItemRepo := mock.NewMockItemRepository(ctrl)
	mockRoosterRepo := mock.NewMockRoosterRepository(ctrl)
	mocKPetRepo := mock.NewMockPetRepository(ctrl)
	mockTxProvider := mock.NewMockTradeTxProvider(ctrl)
	mockCache := mock.NewMockTradeCache(ctrl)

	userService := service.NewUserService(mockUserRepo, mockRoosterRepo, mockItemRepo, mocKPetRepo)
	tradeService := service.NewTradeService(mockCache, userService, mockTxProvider)

	return &TradeServiceTestSuite{
		ctrl:            ctrl,
		ctx:             ctx,
		mockUserRepo:    mockUserRepo,
		mockItemRepo:    mockItemRepo,
		mockRoosterRepo: mockRoosterRepo,
		mockTxProvider:  mockTxProvider,
		mockCache:       mockCache,
		userService:     userService,
		tradeService:    tradeService,
	}
}

func (s *TradeServiceTestSuite) Teardown() {
	s.ctrl.Finish()
}

func TestUpdateUserStatus(t *testing.T) {
	suite := SetupTradeServiceTest(t)
	defer suite.Teardown()

	dto1 := &dto.UpdateUserStatusDTO{
		ID:        uuid.New(),
		Confirmed: true,
		User:      generateSnowflakeLikeID(),
	}

	dto2 := &dto.UpdateUserStatusDTO{
		ID:        uuid.New(),
		Confirmed: true,
		User:      generateSnowflakeLikeID(),
	}

	trade := domain.NewTrade(dto1.ID, dto1.User, dto2.User)
	testCases := []struct {
		name     string
		mockFunc func()
		err      error
		status   *service.UpdateUserStatusWrapper
		dto      *dto.UpdateUserStatusDTO
	}{
		{
			name: "trade not exists",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto1.ID).Return(nil, errors.New("trade not exists"))
			},
			status: nil,
			err:    errors.New("trade not exists"),
			dto:    dto1,
		},
		{
			name: "failed to update cache",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto1.ID).Return(trade, nil)
				suite.mockCache.EXPECT().Update(dto1.ID, gomock.Any()).Return(errors.New("failed to update"))
			},
			status: nil,
			err:    errors.New("failed to update"),
			dto:    dto1,
		},
		{
			name: "sucess",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto1.ID).Return(trade, nil)
				suite.mockCache.EXPECT().Update(dto1.ID, gomock.Any()).Return(nil)
			},
			status: &service.UpdateUserStatusWrapper{
				Trade: trade,
				Done:  false,
			},
			err: nil,
			dto: dto1,
		},
		{
			name: "sucess with trade done",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto2.ID).Return(trade, nil)
				suite.mockCache.EXPECT().Update(dto2.ID, gomock.Any()).Return(nil)
			},
			status: &service.UpdateUserStatusWrapper{
				Trade: trade,
				Done:  true,
			},
			err: nil,
			dto: dto2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()
			status, err := suite.tradeService.UpdateUserStatus(suite.ctx, tc.dto)
			assert.Equal(t, tc.status, status)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestConfirmTrade(t *testing.T) {
	suite := SetupTradeServiceTest(t)
	defer suite.Teardown()

	tradeID, authorID, otherID := uuid.New(), generateSnowflakeLikeID(), generateSnowflakeLikeID()
	trade := domain.NewTrade(tradeID, authorID, otherID)
	seconds := service.COUNTDOWN_SECONDS

	testCases := []struct {
		name          string
		mockFunc      func()
		err           error
		countdownTime int
	}{
		{
			name: "trade not exists",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(tradeID).Return(nil, errors.New("trade not exists"))
			},
			err:           errors.New("trade not exists"),
			countdownTime: 0,
		},
		{
			name: "trade not done",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(tradeID).Return(trade, nil)
			},
			err:           errors.New("trade not done"),
			countdownTime: 0,
		},
		{
			name: "success",
			mockFunc: func() {
				for _, u := range trade.Users {
					u.Confirmed = true
				}
				suite.mockCache.EXPECT().Get(tradeID).Return(trade, nil)
			},
			countdownTime: service.COUNTDOWN_SECONDS,
			err:           nil,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()
			countdownTime, err := suite.tradeService.ConfirmTrade(ctx, tradeID, func(b bool, err error) {})
			assert.Equal(t, tc.countdownTime, countdownTime)
			assert.Equal(t, tc.err, err)
		})
	}

	cancel() // Cancel the context to avoid function callbacks being called twice

	t.Run("context canceled", func(t *testing.T) {

		suite.mockCache.EXPECT().Get(tradeID).Return(trade, nil)

		ctxCancelable, cancel := context.WithCancel(context.Background())

		resultCh := make(chan struct {
			success bool
			err     error
		})

		countdownTime, err := suite.tradeService.ConfirmTrade(ctxCancelable, tradeID, func(success bool, err error) {
			resultCh <- struct {
				success bool
				err     error
			}{success, err}
		})

		assert.Equal(t, seconds, countdownTime)
		assert.Nil(t, err)

		cancel()

		select {
		case result := <-resultCh:
			assert.False(t, result.success)
			assert.Nil(t, result.err)
		case <-time.After(time.Second):
			t.Fatal("Callback was not called after context cancellation")
		}
	})

	t.Run("ticker complete", func(t *testing.T) {
		suite.mockCache.EXPECT().Get(tradeID).Return(trade, nil)
		suite.mockCache.EXPECT().Get(tradeID).Return(trade, nil)
		suite.mockTxProvider.EXPECT().
			Transact(gomock.Any(), gomock.Any()).
			Return(nil)
		suite.mockCache.EXPECT().Delete(tradeID).Return(nil)

		resultCh := make(chan struct {
			success bool
			err     error
		})

		countdownTime, err := suite.tradeService.ConfirmTrade(context.Background(), tradeID, func(success bool, err error) {
			resultCh <- struct {
				success bool
				err     error
			}{success, err}
		})

		assert.Equal(t, seconds, countdownTime)
		assert.Nil(t, err)

		select {
		case result := <-resultCh:
			assert.True(t, result.success)
			assert.Nil(t, result.err)
		case <-time.After(6 * time.Second):
			t.Fatal("Callback was not called after ticker completion")
		}
	})

	t.Run("max of ten roosters", func(t *testing.T) {
		roosterQuantity := 11
		for i := range roosterQuantity {
			rooster := &domain.Rooster{
				ID:     uuid.New(),
				UserID: authorID,
				Type:   i,
			}
			_ = trade.AddItem(authorID, domain.NewTradeItemRooster(rooster))
			suite.mockRoosterRepo.EXPECT().Get(gomock.Any(), rooster.ID).Return(rooster, nil)
			suite.mockRoosterRepo.EXPECT().Delete(gomock.Any(), rooster.ID).Return(nil)
			origin := fmt.Sprintf("Trade with %s", authorID)
			newRooster := domain.NewRooster(otherID, rooster.Type, origin)
			suite.mockRoosterRepo.EXPECT().GetUserRoosterQuantity(gomock.Any(), otherID).Return(i+1, nil)

			suite.mockRoosterRepo.EXPECT().Create(gomock.Any(), newRooster).Return(nil)
		}

		e := errors.New("too many roosters")
		suite.mockCache.EXPECT().Get(tradeID).Return(trade, nil)
		suite.mockCache.EXPECT().Get(tradeID).Return(trade, nil)
		suite.mockTxProvider.EXPECT().
			Transact(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, txFunc func(adapters port.UserTradeTxAdapters, lock func(domain.ID) error) error) error {
				return txFunc(port.UserTradeTxAdapters{
					UserRepository:    suite.mockUserRepo,
					ItemRepository:    suite.mockItemRepo,
					RoosterRepository: suite.mockRoosterRepo,
				}, func(i domain.ID) error {
					return nil
				})
			})

		resultCh := make(chan struct {
			success bool
			err     error
		})

		_, err := suite.tradeService.ConfirmTrade(context.Background(), tradeID, func(success bool, err error) {
			resultCh <- struct {
				success bool
				err     error
			}{success, err}
		})

		assert.Nil(t, err)

		result := <-resultCh
		assert.False(t, result.success)
		assert.Equal(t, e, result.err)
	})

}

func TestUpdateItem(t *testing.T) {
	suite := SetupTradeServiceTest(t)
	defer suite.Teardown()

	authorID, otherID, tradeID, itemID := generateSnowflakeLikeID(), generateSnowflakeLikeID(), uuid.New(), uuid.New()
	dto1 := dto.NewTradeItemDTO(0, tradeID, itemID, authorID, false)
	dto2 := dto.NewTradeItemDTO(0, tradeID, itemID, authorID, true)
	dto3 := dto.NewTradeItemDTO(1, tradeID, itemID, authorID, true)
	dto4 := dto.NewTradeItemDTO(1, tradeID, itemID, authorID, false)

	fakeItem := &domain.Item{
		ID:     itemID,
		UserID: authorID,
		ItemID: 1,
		Type:   domain.NormalType,
	}
	fakeRooster := &domain.Rooster{
		ID:     itemID,
		UserID: authorID,
		Type:   1,
	}
	trade := domain.NewTrade(tradeID, authorID, otherID)
	testCases := []struct {
		name     string
		mockFunc func()
		err      error
		trade    *domain.Trade
		dt       *dto.TradeItemDTO
	}{
		{
			name: "trade not exists",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(tradeID).Return(nil, errors.New("trade not exists"))
			},
			trade: nil,
			err:   errors.New("trade not exists"),
			dt:    dto1,
		},
		{
			name: "trade with invalid item",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto1.ID).Return(trade, nil)
				suite.mockItemRepo.EXPECT().Get(suite.ctx, dto1.ItemID).Return(nil, errors.New("item not found"))
			},
			trade: nil,
			err:   errors.New("item not found"),
			dt:    dto1,
		},
		{
			name: "trade with invalid rooster",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto4.ID).Return(trade, nil)
				suite.mockRoosterRepo.EXPECT().Get(suite.ctx, dto4.ItemID).Return(nil, errors.New("rooster not found"))
			},
			trade: nil,
			err:   errors.New("rooster not found"),
			dt:    dto4,
		},
		{
			name: "sucess dto 1",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto1.ID).Return(trade, nil)
				suite.mockItemRepo.EXPECT().Get(suite.ctx, dto1.ItemID).Return(fakeItem, nil)
				suite.mockCache.EXPECT().Update(dto1.ID, gomock.Any()).Return(nil)
			},
			trade: trade,
			err:   nil,
			dt:    dto1,
		},
		{
			name: "fail dto 1",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto1.ID).Return(trade, nil)
				suite.mockItemRepo.EXPECT().Get(suite.ctx, dto1.ItemID).Return(fakeItem, nil)
			},
			trade: nil,
			err:   errors.New("item quantity exceeded"),
			dt:    dto1,
		},
		{
			name: "sucess dto 2",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto2.ID).Return(trade, nil)
				suite.mockCache.EXPECT().Update(dto2.ID, gomock.Any()).Return(nil)
			},
			trade: trade,
			err:   nil,
			dt:    dto2,
		},
		{
			name: "sucess dto 4",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto4.ID).Return(trade, nil)
				suite.mockRoosterRepo.EXPECT().Get(suite.ctx, dto4.ItemID).Return(fakeRooster, nil)
				suite.mockCache.EXPECT().Update(dto4.ID, gomock.Any()).Return(nil)
			},
			trade: trade,
			err:   nil,
			dt:    dto4,
		},
		{
			name: "fail dto 4",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto4.ID).Return(trade, nil)
				suite.mockRoosterRepo.EXPECT().Get(suite.ctx, dto4.ItemID).Return(fakeRooster, nil)
			},
			trade: nil,
			err:   errors.New("rooster already added"),
			dt:    dto4,
		},
		{
			name: "sucess dto 3",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto3.ID).Return(trade, nil)
				suite.mockCache.EXPECT().Update(dto3.ID, gomock.Any()).Return(nil)
			},
			trade: trade,
			err:   nil,
			dt:    dto3,
		},
		{
			name: "fail dto 3",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto3.ID).Return(trade, nil)
			},
			trade: nil,
			err:   errors.New("rooster not found"),
			dt:    dto3,
		},
		{
			name: "fail dto 2",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto2.ID).Return(trade, nil)
			},
			trade: nil,
			err:   errors.New("item not found"),
			dt:    dto2,
		},
		{
			name: "trade with confirmed user",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(dto1.ID).Return(trade, nil)
				trade.Users[dto1.User].Confirmed = true
			},
			trade: nil,
			err:   errors.New("user already confirmed"),
			dt:    dto1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()
			trade, err := suite.tradeService.UpdateItem(suite.ctx, tc.dt)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.trade, trade)
		})
	}
}

func TestCreateTrade(t *testing.T) {
	suite := SetupTradeServiceTest(t)
	defer suite.Teardown()

	authorID, otherID, tradeID := generateSnowflakeLikeID(), generateSnowflakeLikeID(), uuid.New()

	testCases := []struct {
		name     string
		mockFunc func()
		err      error
		trade    *domain.Trade
	}{
		{
			name: "Trade arleady exists",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(tradeID).Return(&domain.Trade{}, nil)
			},
			trade: nil,
			err:   errors.New("trade already exists"),
		},
		{
			name: "failed to update cache",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(tradeID).Return(nil, nil)
				suite.mockCache.EXPECT().Update(tradeID, gomock.Any()).Return(errors.New("failed to update"))
			},
			trade: nil,
			err:   errors.New("failed to update"),
		},
		{
			name: "sucess",
			mockFunc: func() {
				suite.mockCache.EXPECT().Get(tradeID).Return(nil, nil)
				suite.mockCache.EXPECT().Update(tradeID, gomock.Any()).Return(nil)
			},
			trade: domain.NewTrade(tradeID, authorID, otherID),
			err:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()
			trade, err := suite.tradeService.CreateTrade(suite.ctx, tradeID, authorID, otherID)
			assert.Equal(t, tc.trade, trade)
			assert.Equal(t, tc.err, err)
		})
	}
}
