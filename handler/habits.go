package handler

import (
	"context"
	"fmt"

	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/ZyrnDev/letsgohabits/proto"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
)

func (handler *Handler) HabitsGet(habit *proto.Habit) (*proto.ListHabits, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	habits, err := handler.habitsGRPC.Get(ctx, habit)
	return habits, err
}

func (handler *Handler) HabitsNew(habit *proto.Habit) (*proto.Habit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	habit, err := handler.habitsGRPC.New(ctx, habit)
	return habit, err
}

func (handler *Handler) HabitsDelete(habit *proto.Habit) (*empty.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	empty_, err := handler.habitsGRPC.Delete(ctx, habit)
	return empty_, err
}

func (handler *Handler) HabitUpdate(habit *proto.Habit) (*proto.Habit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	habit, err := handler.habitsGRPC.Update(ctx, habit)
	return habit, err
}

func (handler *Handler) GinFindHabits(c *gin.Context) {
	GinGrpcCall(c, func(habit *database.Habit) (interface{}, error) {
		return handler.HabitsGet(habit.ToProtobuf())
	})
}

func (handler *Handler) GinNewHabit(c *gin.Context) {
	GinGrpcCall(c, func(habit *database.Habit) (interface{}, error) {
		if habit.ID != 0 {
			return nil, fmt.Errorf("habit.ID must be not set or 0")
		}
		if habit.Name == "" {
			return nil, fmt.Errorf("habit.Name must be set and not empty")
		}
		return handler.HabitsNew(habit.ToProtobuf())
	})
}

func (handler *Handler) GinDeleteHabit(c *gin.Context) {
	GinGrpcCall(c, func(habit *database.Habit) (interface{}, error) {
		if habit.ID == 0 {
			return nil, fmt.Errorf("habit.ID must be set and not empty")
		}
		return handler.HabitsDelete(habit.ToProtobuf())
	})
}

func (handler *Handler) GinUpdateHabit(c *gin.Context) {
	GinGrpcCall(c, func(habit *database.Habit) (interface{}, error) {
		if habit.ID == 0 {
			return nil, fmt.Errorf("habit.ID must be set and not empty")
		}
		return handler.HabitUpdate(habit.ToProtobuf())
	})
}
