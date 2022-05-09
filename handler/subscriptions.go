package handler

// import (
// 	"context"
// 	"fmt"

// 	"github.com/ZyrnDev/letsgohabits/database"
// 	"github.com/ZyrnDev/letsgohabits/proto"
// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/protobuf/ptypes/empty"
// )

// func (handler *Handler) SubscriptionsGet(subscription *proto.Subscription) (*proto.ListSubscriptions, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
// 	defer cancel()
// 	subscriptions, err := handler.subscriptionsGRPC.Get(ctx, subscription)
// 	return subscriptions, err
// }

// func (handler *Handler) SubscriptionsNew(subscription *proto.Subscription) (*proto.Subscription, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
// 	defer cancel()
// 	subscription, err := handler.subscriptionsGRPC.New(ctx, subscription)
// 	return subscription, err
// }

// func (handler *Handler) SubscriptionsDelete(subscription *proto.Subscription) (*empty.Empty, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
// 	defer cancel()
// 	empty_, err := handler.subscriptionsGRPC.Delete(ctx, subscription)
// 	return empty_, err
// }

// func (handler *Handler) SubscriptionUpdate(subscription *proto.Subscription) (*proto.Subscription, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
// 	defer cancel()
// 	subscription, err := handler.subscriptionsGRPC.Update(ctx, subscription)
// 	return subscription, err
// }

// func (handler *Handler) GinFindSubscriptions(c *gin.Context) {
// 	GinGrpcCall(c, func(subscription *database.Subscription) (interface{}, error) {
// 		return handler.SubscriptionsGet(subscription.ToProtobuf())
// 	})
// }

// func (handler *Handler) GinNewSubscription(c *gin.Context) {
// 	GinGrpcCall(c, func(subscription *database.Subscription) (interface{}, error) {
// 		if subscription.ID != 0 {
// 			return nil, fmt.Errorf("subscription.ID must be not set or 0")
// 		}
// 		if subscription.UserId == 0 {
// 			return nil, fmt.Errorf("subscription.UserId must be set and not empty")
// 		}
// 		if subscription.HabitId == 0 {
// 			return nil, fmt.Errorf("subscription.HabitId must be set and not empty")
// 		}
// 		if subscription.

// 		return handler.SubscriptionsNew(subscription.ToProtobuf())
// 	})
// }

// func (handler *Handler) GinDeleteSubscription(c *gin.Context) {
// 	GinGrpcCall(c, func(subscription *database.Subscription) (interface{}, error) {
// 		if subscription.ID == 0 {
// 			return nil, fmt.Errorf("subscription.ID must be set and not empty")
// 		}
// 		return handler.SubscriptionsDelete(subscription.ToProtobuf())
// 	})
// }

// func (handler *Handler) GinUpdateSubscription(c *gin.Context) {
// 	GinGrpcCall(c, func(subscription *database.Subscription) (interface{}, error) {
// 		if subscription.ID == 0 {
// 			return nil, fmt.Errorf("subscription.ID must be set and not empty")
// 		}
// 		return handler.SubscriptionUpdate(subscription.ToProtobuf())
// 	})
// }
