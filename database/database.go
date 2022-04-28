package database

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ZyrnDev/letsgohabits/proto"
)

type Database = *gorm.DB
type Config = gorm.Config

type Model struct {
	ID        uint `gorm:"primarykey form:"id" json:"id""`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	Model
	Nickname string  `gorm:"unique;index" form:"nickname" json:"nickname"`
	Habits   []Habit `gorm:"foreignKey:AuthorId"`
}

type Habit struct {
	Model
	Name        string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
	AuthorId    uint   `form:"authorId" json:"authorId"`
	// Events []Event
}

func New(connectionString string, conf *Config) (Database, error) {
	db, err := gorm.Open(sqlite.Open(connectionString), conf)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{}, &Habit{})

	return db, nil
}

func (user *User) ToProtobuf() *proto.User {
	return &proto.User{
		Id:        uint64(user.ID),
		Name:      user.Nickname,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		// DeletedAt: timestamppb.New(user.DeletedAt),
	}
}

func (user *User) FromProtobuf(in *proto.User) {
	user.ID = uint(in.Id)
	user.Nickname = in.Name
	user.CreatedAt = in.CreatedAt.AsTime()
	user.UpdatedAt = in.UpdatedAt.AsTime()
}

func (habit *Habit) ToProtobuf() *proto.Habit {
	return &proto.Habit{
		Id:          uint64(habit.ID),
		Name:        habit.Name,
		Description: habit.Description,
		AuthorId:    uint64(habit.AuthorId),
		CreatedAt:   timestamppb.New(habit.CreatedAt),
		UpdatedAt:   timestamppb.New(habit.UpdatedAt),
		// DeletedAt: timestamppb.New(habit.DeletedAt),
	}
}

func (habit *Habit) FromProtobuf(in *proto.Habit) {
	habit.ID = uint(in.Id)
	habit.Name = in.Name
	habit.Description = in.Description
	habit.AuthorId = uint(in.AuthorId)
	habit.CreatedAt = in.CreatedAt.AsTime()
	habit.UpdatedAt = in.UpdatedAt.AsTime()
}
