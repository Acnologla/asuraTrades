package service_test

import (
	"testing"
	"time"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/acnologla/asuraTrades/internal/core/port/mock"
	"github.com/acnologla/asuraTrades/internal/core/service"
	"github.com/brianvoe/gofakeit/v6"
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

	userTrade, _ := domain.NewUserTrade(dto.AuthorID, dto.OtherID, dto.TradeID)

	service := service.NewUserTokenService(mockTokenProvider, mockUserRepo, mockItemRepo, mockRoosterRepo)

	fakeToken := gofakeit.UUID()

	mockUserRepo.EXPECT().Get(gomock.Any(), gomock.AnyOf(somerandomID, somerandomID2)).DoAndReturn(func(x any, id domain.ID) (*domain.User, error) {
		return &domain.User{
			ID: id,
		}, nil
	}).Times(2)

	mockTokenProvider.EXPECT().GenerateToken(userTrade, gomock.Any()).Return(fakeToken, nil)

	token, err := service.CreateToken(ctx, dto)

	assert.Equal(t, fakeToken, token)
	assert.Nil(t, err)
}
