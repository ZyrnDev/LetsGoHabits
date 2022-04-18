package server

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/ZyrnDev/letsgohabits/grpc"
	"github.com/ZyrnDev/letsgohabits/nats"
	"github.com/ZyrnDev/letsgohabits/proto"
	"github.com/go-co-op/gocron"
	googleGrpc "google.golang.org/grpc"
	protobuf "google.golang.org/protobuf/proto"
)

func New(natsConnStr string, db database.Database, shutdownRequested chan bool) chan bool {
	done := make(chan bool)

	go func() {
		nc := nats.NatsConnection(natsConnStr)
		defer nc.Close()

		log.Printf("Server connected to %s", natsConnStr)

		scheduler := gocron.NewScheduler(time.Local)

		// scheduler.Every(1).Seconds().Do(func() {
		// 	nc.Publish("print", []byte(fmt.Sprintf("The time is: %s", time.Now())))
		// })

		scheduler.CronWithSeconds("*/2 * * * * *").Do(func() {
			habit := database.Habit{
				Name:        fmt.Sprintf("My New Habit %d", rand.Intn(256)),
				Description: "a cool as heck habit",
				Author: database.User{
					Nickname: "ZyrnDev (Mitch)",
				},
			}

			res := db.Create(&habit)

			if res.Error != nil {
				panic(res.Error)
			} else if res.RowsAffected != 1 {
				panic("Rows affected is not 1")
			}

			t := &proto.Test{
				Name: fmt.Sprintf("Created %d", habit.ID),
				Id:   uint64(habit.ID),
			}

			msg, err := protobuf.Marshal(t)
			if err != nil {
				panic(err)
			}

			nc.Publish("test", msg)
			// nc.Publish("print", []byte(fmt.Sprintf("The time is: %s", time.Now())))
		})

		go test()

		scheduler.StartAsync()

		<-shutdownRequested
		done <- true
	}()

	return done
}

func test() {
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []googleGrpc.ServerOption
	grpcServer := googleGrpc.NewServer(opts...)
	proto.RegisterToolsServer(grpcServer, &grpc.ToolsServer{})
	grpcServer.Serve(lis)

}
