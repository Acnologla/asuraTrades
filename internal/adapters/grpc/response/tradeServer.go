package response

import (
	"github.com/acnologla/asuraTrades/internal/adapters/grpc/proto"
	"github.com/acnologla/asuraTrades/internal/core/domain"
)

func NewFinishTradeResponse(err error) *proto.FinishTradeResponse {
	if err == nil {
		return &proto.FinishTradeResponse{
			Ok: true,
		}
	}

	errStr := err.Error()
	return &proto.FinishTradeResponse{
		Ok:    false,
		Error: &errStr,
	}
}

func NewUserProfileResponse(profile *domain.UserProfile) *proto.GetUserProfileResponse {
	roosters := make([]*proto.Rooster, len(profile.Roosters))
	items := make([]*proto.Item, len(profile.Items))

	for i, rooster := range profile.Roosters {
		roosters[i] = &proto.Rooster{
			Id:     rooster.ID.String(),
			Origin: rooster.Origin,
			Type:   int32(rooster.Type),
			UserId: uint64(rooster.UserID),
		}
	}

	for i, item := range profile.Items {
		items[i] = &proto.Item{
			Id:       item.ID.String(),
			Type:     int32(item.Type),
			ItemId:   int32(item.ItemID),
			UserId:   uint64(item.UserID),
			Quantity: int32(item.Quantity),
		}
	}

	return &proto.GetUserProfileResponse{
		Id:       uint64(profile.ID),
		Roosters: roosters,
		Items:    items,
	}
}
