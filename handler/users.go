package handler

import (
	"context"
	"fmt"

	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/ZyrnDev/letsgohabits/proto"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
)

// TODO(Mitch): Make this more generic using a wrapper function
func (handler *Handler) UsersGet(user *proto.User) (*proto.ListUsers, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	users, err := handler.usersGRPC.Get(ctx, user)
	return users, err
}

func (handler *Handler) UsersNew(user *proto.User) (*proto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	user, err := handler.usersGRPC.New(ctx, user)
	return user, err
}

func (handler *Handler) UsersDelete(user *proto.User) (*empty.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	empty_, err := handler.usersGRPC.Delete(ctx, user)
	return empty_, err
}

func (handler *Handler) UserUpdate(user *proto.User) (*proto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	user, err := handler.usersGRPC.Update(ctx, user)
	return user, err
}

func (handler *Handler) GinFindUsers(c *gin.Context) {
	GinGrpcCall(c, func(user *database.User) (interface{}, error) {
		return handler.UsersGet(user.ToProtobuf())
	})
}

func (handler *Handler) GinNewUser(c *gin.Context) {
	GinGrpcCall(c, func(user *database.User) (interface{}, error) {
		if user.ID != 0 {
			return nil, fmt.Errorf("user.ID must be not set or 0")
		}
		if user.Nickname == "" {
			return nil, fmt.Errorf("user.Nickname must be set and not empty")
		}
		return handler.UsersNew(user.ToProtobuf())
	})
}

func (handler *Handler) GinDeleteUser(c *gin.Context) {
	GinGrpcCall(c, func(user *database.User) (interface{}, error) {
		if user.ID == 0 {
			return nil, fmt.Errorf("user.ID must be set and not empty")
		}
		return handler.UsersDelete(user.ToProtobuf())
	})
}

func (handler *Handler) GinUpdateUser(c *gin.Context) {
	GinGrpcCall(c, func(user *database.User) (interface{}, error) {
		if user.ID == 0 {
			return nil, fmt.Errorf("user.ID must be set and not empty")
		}
		return handler.UserUpdate(user.ToProtobuf())
	})
}
