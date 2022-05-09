package grpc

import (
	"context"

	// "github.com/ZyrnDev/letsgohabits/mounts/proto/proto"
	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/ZyrnDev/letsgohabits/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type UsersServer struct {
	proto.UnimplementedUsersServer
	Database database.Database
}

func (s *UsersServer) New(ctx context.Context, in *proto.User) (*proto.User, error) {
	return New[proto.User, database.User, *database.User](ctx, s.Database, in)
}

func (s *UsersServer) Delete(ctx context.Context, in *proto.User) (*empty.Empty, error) {
	return Delete[proto.User, database.User, *database.User](ctx, s.Database, in)
}

func (s *UsersServer) Get(ctx context.Context, in *proto.User) (*proto.ListUsers, error) {
	results, err := Get[proto.User, database.User, *database.User](ctx, s.Database, in)

	if err != nil {
		return nil, err
	}

	return &proto.ListUsers{Users: results}, nil
}

func (s *UsersServer) Update(ctx context.Context, in *proto.User) (*proto.User, error) {
	return Update[proto.User, database.User, *database.User](ctx, s.Database, in)
}
