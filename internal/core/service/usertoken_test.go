package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/acnologla/asuraTrades/internal/core/port/mock"
	"github.com/acnologla/asuraTrades/internal/core/service"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func generateSnowflakeLikeID() domain.ID {
	timestamp := uint64(time.Now().UnixNano() / int64(time.Millisecond))
	machineID := uint64(gofakeit.Int64())
	sequence := uint64(gofakeit.Int64())
	snowflakeID := (timestamp << 22) | (machineID << 12) | sequence
	return domain.ID(snowflakeID)
}

func TestGetTradeTokenResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := t.Context()
	defer ctrl.Finish()

	mockTokenProvider := mock.NewMockTokenProvider(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	mockItemRepo := mock.NewMockItemRepository(ctrl)
	mockRoosterRepo := mock.NewMockRoosterRepository(ctrl)

	somerandomID := generateSnowflakeLikeID()
	somerandomID2 := generateSnowflakeLikeID()

	tradeID := gofakeit.UUID()
	trade, _ := domain.NewUserTrade(somerandomID.String(), somerandomID2.String(), tradeID)

	fakeRooster := &domain.Rooster{
		ID: uuid.New(),
	}

	fakeItem := &domain.Item{
		ID: uuid.New(),
	}

	response := &service.GetTradeTokenResponseWrapper{
		UserTrade: trade,
		UserProfile: &domain.UserProfile{
			User: &domain.User{
				ID: somerandomID,
			},
			Roosters: []*domain.Rooster{fakeRooster},
			Items:    []*domain.Item{fakeItem},
		},
	}

	fakeToken := gofakeit.UUID()
	testCases := []struct {
		name     string
		mockFunc func()
		err      error
		res      *service.GetTradeTokenResponseWrapper
	}{
		{
			name: "Invalid Token",
			mockFunc: func() {
				mockTokenProvider.EXPECT().ValidateToken(gomock.Any()).Return(nil, errors.New("Invalid token"))
			},
			res: nil,
			err: errors.New("Invalid token"),
		},
		{
			name: "Invalid User",
			mockFunc: func() {
				mockTokenProvider.EXPECT().ValidateToken(gomock.Any()).Return(trade, nil)

				mockUserRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("Invalid user"))
			},
			res: nil,
			err: errors.New("Invalid user"),
		},
		{
			name: "Invalid items",
			mockFunc: func() {
				mockTokenProvider.EXPECT().ValidateToken(gomock.Any()).Return(trade, nil)
				mockUserRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&domain.User{
					ID: somerandomID,
				}, nil)
				mockItemRepo.EXPECT().GetUserItems(gomock.Any(), somerandomID).Return([]*domain.Item{}, errors.New("invalid items"))
			},
			res: nil,
			err: errors.New("invalid items"),
		},
		{
			name: "Invalid roosters",
			mockFunc: func() {
				mockTokenProvider.EXPECT().ValidateToken(gomock.Any()).Return(trade, nil)
				mockUserRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&domain.User{
					ID: somerandomID,
				}, nil)
				mockItemRepo.EXPECT().GetUserItems(gomock.Any(), somerandomID).Return([]*domain.Item{}, nil)
				mockRoosterRepo.EXPECT().GetUserRoosters(gomock.Any(), somerandomID).Return([]*domain.Rooster{}, errors.New("invalid roosters"))

			},
			res: nil,
			err: errors.New("invalid roosters"),
		},
		{

			name: "Success",
			mockFunc: func() {
				mockTokenProvider.EXPECT().ValidateToken(gomock.Any()).Return(trade, nil)
				mockUserRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&domain.User{
					ID: somerandomID,
				}, nil)

				mockItemRepo.EXPECT().GetUserItems(gomock.Any(), somerandomID).Return(response.UserProfile.Items, nil)
				mockRoosterRepo.EXPECT().GetUserRoosters(gomock.Any(), somerandomID).Return(response.UserProfile.Roosters, nil)

			},
			res: response,
			err: nil,
		},
	}

	service := service.NewUserTokenService(mockTokenProvider, mockUserRepo, mockItemRepo, mockRoosterRepo)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()
			res, err := service.GetTradeTokenResponse(ctx, fakeToken)
			assert.Equal(t, tc.res, res)
			assert.Equal(t, tc.err, err)
		})

	}
}

func TestCreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := t.Context()
	defer ctrl.Finish()

	mockTokenProvider := mock.NewMockTokenProvider(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	mockItemRepo := mock.NewMockItemRepository(ctrl)
	mockRoosterRepo := mock.NewMockRoosterRepository(ctrl)

	somerandomID := generateSnowflakeLikeID()
	somerandomID2 := generateSnowflakeLikeID()
	dto := &dto.GenerateUserTokenDTO{
		AuthorID: somerandomID.String(),
		OtherID:  somerandomID2.String(),
		TradeID:  gofakeit.UUID(),
	}
	fakeToken := gofakeit.UUID()
	testCases := []struct {
		name     string
		mockFunc func()
		err      error
		token    string
	}{
		{
			name: "Invalid Token",
			mockFunc: func() {
				mockUserRepo.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(func(x any, id domain.ID) (*domain.User, error) {
					return &domain.User{}, nil
				}).Times(2)

				mockTokenProvider.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).Return("", errors.New("Invalid token"))
			},
			token: "",
			err:   errors.New("Invalid token"),
		},
		{
			name: "Invalid User",
			mockFunc: func() {
				mockUserRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("Invalid user"))
			},
			token: "",
			err:   errors.New("Invalid user"),
		},
		{

			name: "Success",
			mockFunc: func() {
				mockUserRepo.EXPECT().Get(gomock.Any(), gomock.AnyOf(somerandomID, somerandomID2)).DoAndReturn(func(x any, id domain.ID) (*domain.User, error) {
					return &domain.User{
						ID: id,
					}, nil
				}).Times(2)
				userTrade, _ := domain.NewUserTrade(dto.AuthorID, dto.OtherID, dto.TradeID)

				mockTokenProvider.EXPECT().GenerateToken(userTrade, gomock.Any()).Return(fakeToken, nil)
			},
			token: fakeToken,
			err:   nil,
		},
	}

	service := service.NewUserTokenService(mockTokenProvider, mockUserRepo, mockItemRepo, mockRoosterRepo)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()
			token, err := service.CreateToken(ctx, dto)
			assert.Equal(t, tc.token, token)
			assert.Equal(t, tc.err, err)
		})

	}
}
