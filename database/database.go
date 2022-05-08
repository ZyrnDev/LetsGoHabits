package database

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/SherClockHolmes/webpush-go"
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
	Nickname      string         `gorm:"unique;index" form:"nickname" json:"nickname"`
	Habits        []Habit        `gorm:"foreignKey:AuthorId"`
	Subscriptions []Subscription `gorm:"foreignKey:UserId"`
}

type Habit struct {
	Model
	Name          string         `form:"name" json:"name"`
	Description   string         `form:"description" json:"description"`
	AuthorId      uint           `form:"authorId" json:"authorId"`
	Subscriptions []Subscription `gorm:"foreignKey:HabitId"`
	// Events []Event
}

type Subscription struct {
	Model
	UserId   uint   `form:"userId" json:"userId"`
	HabitId  uint   `form:"habitId" json:"habitId"`
	Endpoint string `form:"endpoint" json:"endpoint"`
	Keys     `form:"keys" json:"keys"`
}

type Keys struct {
	Auth   string `json:"auth"`
	P256dh string `json:"p256dh"`
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

func (subscription *Subscription) ToProtobuf() *proto.Subscription {
	return &proto.Subscription{
		Id:        uint64(subscription.ID),
		UserId:    uint64(subscription.UserId),
		HabitId:   uint64(subscription.HabitId),
		Endpoint:  subscription.Endpoint,
		Keys:      subscription.Keys.ToProtobuf(),
		CreatedAt: timestamppb.New(subscription.CreatedAt),
		UpdatedAt: timestamppb.New(subscription.UpdatedAt),
		// DeletedAt: timestamppb.New(subscription.DeletedAt),
	}
}

func (subscription *Subscription) FromProtobuf(in *proto.Subscription) {
	subscription.ID = uint(in.Id)
	subscription.UserId = uint(in.UserId)
	subscription.HabitId = uint(in.HabitId)
	subscription.Endpoint = in.Endpoint
	// subscription.Keys.FromProtobuf(in.Keys)
	subscription.CreatedAt = in.CreatedAt.AsTime()
	subscription.UpdatedAt = in.UpdatedAt.AsTime()
}

func (keys *Keys) ToProtobuf() *proto.Subscription_Keys {
	return &proto.Subscription_Keys{
		Auth:   keys.Auth,
		P256Dh: keys.P256dh,
	}
}

func (keys *Keys) FromProtobuf(in *proto.Subscription_Keys) {
	keys.Auth = in.Auth
	keys.P256dh = in.P256Dh
}

func (keys *Keys) ToWebpushKeys() webpush.Keys {
	return webpush.Keys{
		Auth:   keys.Auth,
		P256dh: keys.P256dh,
	}
}

func (sub *Subscription) ToWebpushSubScription() webpush.Subscription {
	return webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys:     sub.Keys.ToWebpushKeys(),
	}
}
