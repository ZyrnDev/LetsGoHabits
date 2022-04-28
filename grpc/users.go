package grpc

import (
	"context"

	// "github.com/ZyrnDev/letsgohabits/mounts/proto/proto"
	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/ZyrnDev/letsgohabits/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersServer struct {
	proto.UnimplementedUsersServer
	Database database.Database
}

func (s *UsersServer) New(ctx context.Context, in *proto.User) (*proto.User, error) {
	var user database.User
	user.FromProtobuf(in)

	res := s.Database.Create(&user)

	if res.Error != nil || res.RowsAffected != 1 {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", res.Error)
	}

	return user.ToProtobuf(), nil
}

func (s *UsersServer) Delete(ctx context.Context, in *proto.User) (*empty.Empty, error) {
	var user database.User
	user.FromProtobuf(in)

	result := s.Database.Unscoped().Delete(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %s", result.Error)
	}

	return &empty.Empty{}, nil
}

func (s *UsersServer) Get(ctx context.Context, in *proto.User) (*proto.ListUsers, error) {
	var user database.User
	user.FromProtobuf(in)

	var users []database.User

	res := s.Database.Find(&users, user)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %s", res.Error)
	}

	usersProto := make([]*proto.User, len(users))
	for i, user := range users {
		usersProto[i] = user.ToProtobuf()
	}

	return &proto.ListUsers{
		Users: usersProto,
	}, nil

}

func (s *UsersServer) Update(ctx context.Context, in *proto.User) (*proto.User, error) {
	var user database.User
	user.FromProtobuf(in)

	var updatedUser database.User
	if res := s.Database.First(&updatedUser, user.ID); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user to update: %s", res.Error)
	}

	if res := s.Database.Model(&updatedUser).Updates(user); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", res.Error)
	}

	if res := s.Database.First(&updatedUser, user.ID); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user after update: %s", res.Error)
	}

	return updatedUser.ToProtobuf(), nil
}
