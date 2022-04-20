package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database = *gorm.DB
type Config = gorm.Config

type User struct {
	gorm.Model
	Nickname string
}

type Habit struct {
	gorm.Model
	Name        string
	Description string
	AuthorID    uint
	Author      User
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
