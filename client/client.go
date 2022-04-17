package client

import (
	"log"

	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/ZyrnDev/letsgohabits/nats"
	"github.com/ZyrnDev/letsgohabits/proto"
	protobuf "google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

func New(natsConnStr string, db database.Database, shutdownRequested chan bool) chan bool {
	done := make(chan bool)

	go func() {
		nc := nats.NatsConnection(natsConnStr)
		defer nc.Close()

		log.Printf("Client connected to %s", natsConnStr)

		nc.Subscribe("print", func(msg nats.NatsMsg) {
			log.Printf("Received: %s", string(msg.Data))
		})

		nc.Subscribe("test", func(msg nats.NatsMsg) {
			t := &proto.Test{}
			err := protobuf.Unmarshal(msg.Data, t)
			if err != nil {
				panic(err)
			}

			var habit database.Habit
			res := db.Preload("Author").First(&habit, database.Habit{Model: gorm.Model{ID: uint(t.Id)}})
			if res.Error != nil {
				panic(res.Error)
			}

			log.Printf("%s just send a message (%s) to add a new habit '%s'", habit.Author.Nickname, t.Name, habit.Name)
			// log.Printf("Info %+v", habit)
		})

		<-shutdownRequested
		done <- true
	}()

	return done
}
