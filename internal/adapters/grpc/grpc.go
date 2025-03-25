package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/acnologla/asuraTrades/internal/adapters/grpc/proto"
	"google.golang.org/grpc"
)

//go:generate protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. ./proto/trade.proto

type server struct{}

func (s *server) GetUserProfile(ctx context.Context, in *proto.GetUserProfileRequest) (*proto.GetUserProfileResponse, error) {
	return &proto.GetUserProfileResponse{}, nil
}

func (s *server) FinishTrade(ctx context.Context, in *proto.FinishTradeRequest) (*proto.FinishTradeResponse, error) {
	return &proto.FinishTradeResponse{}, nil
}

func NewGrpcServer(port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	s := grpc.NewServer()
	proto.RegisterTradeServer(s, &server{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
