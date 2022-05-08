package grpc

import (
	"context"

	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/ZyrnDev/letsgohabits/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubscriptionsServer struct {
	proto.UnimplementedSubscriptionsServer
	Database database.Database
}

func (s *SubscriptionsServer) New(ctx context.Context, in *proto.Subscription) (*proto.Subscription, error) {
	var subscription database.Subscription
	subscription.FromProtobuf(in)

	res := s.Database.Create(&subscription)

	if res.Error != nil || res.RowsAffected != 1 {
		return nil, status.Errorf(codes.Internal, "failed to create subscription: %s", res.Error)
	}

	return subscription.ToProtobuf(), nil
}

func (s *SubscriptionsServer) Delete(ctx context.Context, in *proto.Subscription) (*empty.Empty, error) {
	var subscription database.Subscription
	subscription.FromProtobuf(in)

	result := s.Database.Unscoped().Delete(&subscription)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete subscription: %s", result.Error)
	}

	return &empty.Empty{}, nil
}

func (s *SubscriptionsServer) Get(ctx context.Context, in *proto.Subscription) (*proto.ListSubscriptions, error) {
	var subscription database.Subscription
	subscription.FromProtobuf(in)

	var subscriptions []database.Subscription

	res := s.Database.Find(&subscriptions, subscription)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get subscription: %s", res.Error)
	}

	subscriptionsProto := make([]*proto.Subscription, len(subscriptions))
	for i, subscription := range subscriptions {
		subscriptionsProto[i] = subscription.ToProtobuf()
	}

	return &proto.ListSubscriptions{
		Subscriptions: subscriptionsProto,
	}, nil

}

func (s *SubscriptionsServer) Update(ctx context.Context, in *proto.Subscription) (*proto.Subscription, error) {
	var subscription database.Subscription
	subscription.FromProtobuf(in)

	var updatedSubscription database.Subscription
	if res := s.Database.First(&updatedSubscription, subscription.ID); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to find subscription to update: %s", res.Error)
	}

	if res := s.Database.Model(&updatedSubscription).Updates(subscription); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to update subscription: %s", res.Error)
	}

	if res := s.Database.First(&updatedSubscription, subscription.ID); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to find subscription after update: %s", res.Error)
	}

	return updatedSubscription.ToProtobuf(), nil
}
