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

type HabitsServer struct {
	proto.UnimplementedHabitsServer
	Database database.Database
}

func (s *HabitsServer) New(ctx context.Context, in *proto.Habit) (*proto.Habit, error) {
	var habit database.Habit
	habit.FromProtobuf(in)

	res := s.Database.Create(&habit)

	if res.Error != nil || res.RowsAffected != 1 {
		return nil, status.Errorf(codes.Internal, "failed to create habit: %s", res.Error)
	}

	return habit.ToProtobuf(), nil
}

func (s *HabitsServer) Delete(ctx context.Context, in *proto.Habit) (*empty.Empty, error) {
	var habit database.Habit
	habit.FromProtobuf(in)

	result := s.Database.Unscoped().Delete(&habit)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete habit: %s", result.Error)
	}

	return &empty.Empty{}, nil
}

func (s *HabitsServer) Get(ctx context.Context, in *proto.Habit) (*proto.ListHabits, error) {
	var habit database.Habit
	habit.FromProtobuf(in)

	var habits []database.Habit

	res := s.Database.Find(&habits, habit)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get habit: %s", res.Error)
	}

	habitsProto := make([]*proto.Habit, len(habits))
	for i, habit := range habits {
		habitsProto[i] = habit.ToProtobuf()
	}

	return &proto.ListHabits{
		Habits: habitsProto,
	}, nil

}

func (s *HabitsServer) Update(ctx context.Context, in *proto.Habit) (*proto.Habit, error) {
	var habit database.Habit
	habit.FromProtobuf(in)

	var updatedHabit database.Habit
	if res := s.Database.First(&updatedHabit, habit.ID); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to find habit to update: %s", res.Error)
	}

	if res := s.Database.Model(&updatedHabit).Updates(habit); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to update habit: %s", res.Error)
	}

	if res := s.Database.First(&updatedHabit, habit.ID); res.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to find habit after update: %s", res.Error)
	}

	return updatedHabit.ToProtobuf(), nil
}
