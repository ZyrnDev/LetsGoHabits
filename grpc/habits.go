package grpc

import (
	"context"

	// "github.com/ZyrnDev/letsgohabits/mounts/proto/proto"
	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/ZyrnDev/letsgohabits/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type HabitsServer struct {
	proto.UnimplementedHabitsServer
	Database database.Database
}

func (s *HabitsServer) New(ctx context.Context, in *proto.Habit) (*proto.Habit, error) {
	return New[proto.Habit, database.Habit, *database.Habit](ctx, s.Database, in)
}

func (s *HabitsServer) Delete(ctx context.Context, in *proto.Habit) (*empty.Empty, error) {
	return Delete[proto.Habit, database.Habit, *database.Habit](ctx, s.Database, in)
}

func (s *HabitsServer) Get(ctx context.Context, in *proto.Habit) (*proto.ListHabits, error) {
	results, err := Get[proto.Habit, database.Habit, *database.Habit](ctx, s.Database, in)

	if err != nil {
		return nil, err
	}

	return &proto.ListHabits{Habits: results}, nil

}

func (s *HabitsServer) Update(ctx context.Context, in *proto.Habit) (*proto.Habit, error) {
	return Update[proto.Habit, database.Habit, *database.Habit](ctx, s.Database, in)
}
