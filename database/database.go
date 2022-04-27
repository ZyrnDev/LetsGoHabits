package database

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ZyrnDev/letsgohabits/proto"
)

type Database = *gorm.DB
type Config = gorm.Config

type User struct {
	gorm.Model
	Nickname string  `gorm:"unique;index"`
	Habits   []Habit `gorm:"foreignKey:AuthorId"`
}

type Habit struct {
	gorm.Model
	Name        string
	Description string
	AuthorId    uint
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
