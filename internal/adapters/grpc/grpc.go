package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/acnologla/asuraTrades/internal/adapters/config"
	"github.com/acnologla/asuraTrades/internal/adapters/grpc/proto"
	"github.com/acnologla/asuraTrades/internal/adapters/grpc/response"
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/dto"
	"github.com/acnologla/asuraTrades/internal/core/service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

//go:generate protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. ./proto/trade.proto

type server struct {
	userService  *service.UserService
	tradeService *service.TradeService
}

func (s *server) GetUserProfile(ctx context.Context, in *proto.GetUserProfileRequest) (*proto.GetUserProfileResponse, error) {
	profile, err := s.userService.GetUserProfile(ctx, domain.ID(in.Id))
	if err != nil {
		return nil, err
	}
	return response.NewUserProfileResponse(profile), nil
}

func (s *server) addItemsToTrade(ctx context.Context, tradeID uuid.UUID, items []*proto.TradeItem) error {
	for _, item := range items {
		itemID, err := uuid.Parse(item.Id)
		if err != nil {
			return errors.New("invalid item id")
		}

		data := dto.NewTradeItemDTO(int(item.Type), tradeID, itemID, domain.ID(item.UserId), false)
		_, err = s.tradeService.UpdateItem(ctx, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server) confirmUsersStatus(ctx context.Context, in *proto.FinishTradeRequest, tradeID uuid.UUID) {
	for _, userID := range []uint64{in.AuthorId, in.OtherId} {
		s.tradeService.UpdateUserStatus(ctx, dto.NewUpdateUserStatusDTO(tradeID, true, domain.ID(userID)))
	}
}

func (s *server) FinishTrade(ctx context.Context, in *proto.FinishTradeRequest) (*proto.FinishTradeResponse, error) {
	tradeID := uuid.New()
	_, err := s.tradeService.CreateTrade(ctx, tradeID, domain.ID(in.AuthorId), domain.ID(in.OtherId))

	if err != nil {
		return nil, err
	}

	if err := s.addItemsToTrade(ctx, tradeID, in.Items); err != nil {
		return nil, err
	}

	s.confirmUsersStatus(ctx, in, tradeID)
	err = s.tradeService.FinishTrade(ctx, tradeID)

	return response.NewFinishTradeResponse(err), nil
}

func NewGrpcServer(grpcConfig config.GrpcConfig, userService *service.UserService, tradeService *service.TradeService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", grpcConfig.Port, err)
	}

	interceptor := NewAuthInterceptor(grpcConfig.Token)
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryInterceptor))

	proto.RegisterTradeServer(s, &server{
		userService:  userService,
		tradeService: tradeService,
	})

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
