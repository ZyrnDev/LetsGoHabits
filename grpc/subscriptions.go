package grpc

import (
	"context"

	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/ZyrnDev/letsgohabits/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type SubscriptionsServer struct {
	proto.UnimplementedSubscriptionsServer
	Database database.Database
}

func (s *SubscriptionsServer) New(ctx context.Context, in *proto.Subscription) (*proto.Subscription, error) {
	return New[proto.Subscription, database.Subscription, *database.Subscription](ctx, s.Database, in)
}

func (s *SubscriptionsServer) Delete(ctx context.Context, in *proto.Subscription) (*empty.Empty, error) {
	return Delete[proto.Subscription, database.Subscription, *database.Subscription](ctx, s.Database, in)
}

func (s *SubscriptionsServer) Get(ctx context.Context, in *proto.Subscription) (*proto.ListSubscriptions, error) {
	results, err := Get[proto.Subscription, database.Subscription, *database.Subscription](ctx, s.Database, in)

	if err != nil {
		return nil, err
	}

	return &proto.ListSubscriptions{Subscriptions: results}, nil
}

func (s *SubscriptionsServer) Update(ctx context.Context, in *proto.Subscription) (*proto.Subscription, error) {
	return Update[proto.Subscription, database.Subscription, *database.Subscription](ctx, s.Database, in)
}
