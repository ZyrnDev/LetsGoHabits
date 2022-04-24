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
	user := &database.User{
		Nickname: in.Name,
	}

	res := s.Database.Create(&user)

	if res.Error != nil || res.RowsAffected != 1 {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", res.Error)
	}

	return &proto.User{
		Id:   uint64(user.ID),
		Name: user.Nickname,
	}, nil
}

func (s *UsersServer) Delete(ctx context.Context, in *proto.User) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func (s *UsersServer) Get(ctx context.Context, in *proto.User) (*proto.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}

func (s *UsersServer) Update(ctx context.Context, in *proto.UpdateUserRequest) (*proto.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
