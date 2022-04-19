package server

import (
	"time"

	"github.com/ZyrnDev/letsgohabits/config"
	"github.com/ZyrnDev/letsgohabits/nats"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

type DatabaseConfig struct {
	ConnectionString string `mapstructure:"connection_string"`
}

type NatsConfig struct {
	ConnectionString string `mapstructure:"connection_string"`
}

type ServerConfig struct {
	DatabaseConfig `mapstructure:"database"`
	NatsConfig     `mapstructure:"nats"`
}

func New(args ...string) {
	conf, err := config.New[ServerConfig]("server.toml")
	if err != nil {
		log.Fatal().Msgf("Failed to load config: %s", err)
	} else {
		log.Info().Msgf("Loaded Config: %+v", conf)
	}

	// db := database.New(conf.DatabaseConnectionString, &database.Config{
	// 	// Logger: logger.Default.LogMode(logger.Info), // Verbose Logging
	// })

	log.Info().Msgf("Starting server %+v", args)

	nc, err := nats.Connect(conf.NatsConfig.ConnectionString)
	if err != nil {
		log.Fatal().Msgf("Failed to connect to nats: %s", err)
	}
	// defer nc.Close() // TODO: This is causing issue as it closes too early

	log.Info().Msgf("Connected to nats: %s", conf.NatsConfig.ConnectionString)

	scheduler := gocron.NewScheduler(time.Local)
	scheduler.Every(1).Second().Do(func() {
		nc.Publish("print", []byte("Hello World"))
	})
	scheduler.StartAsync()
}

// func New(natsConnStr string, db database.Database, shutdownRequested chan bool) chan bool {
// 	done := make(chan bool)

// 	go func() {
// 		nc := nats.NatsConnection(natsConnStr)
// 		defer nc.Close()

// 		log.Printf("Server connected to %s", natsConnStr)

// 		scheduler := gocron.NewScheduler(time.Local)

// 		// scheduler.Every(1).Seconds().Do(func() {
// 		// 	nc.Publish("print", []byte(fmt.Sprintf("The time is: %s", time.Now())))
// 		// })

// 		scheduler.CronWithSeconds("*/2 * * * * *").Do(func() {
// 			habit := database.Habit{
// 				Name:        fmt.Sprintf("My New Habit %d", rand.Intn(256)),
// 				Description: "a cool as heck habit",
// 				Author: database.User{
// 					Nickname: "ZyrnDev (Mitch)",
// 				},
// 			}

// 			res := db.Create(&habit)

// 			if res.Error != nil {
// 				panic(res.Error)
// 			} else if res.RowsAffected != 1 {
// 				panic("Rows affected is not 1")
// 			}

// 			t := &proto.Test{
// 				Name: fmt.Sprintf("Created %d", habit.ID),
// 				Id:   uint64(habit.ID),
// 			}

// 			msg, err := protobuf.Marshal(t)
// 			if err != nil {
// 				panic(err)
// 			}

// 			nc.Publish("test", msg)
// 			// nc.Publish("print", []byte(fmt.Sprintf("The time is: %s", time.Now())))
// 		})

// 		go test()

// 		scheduler.StartAsync()

// 		<-shutdownRequested
// 		done <- true
// 	}()

// 	return done
// }

// func test() {
// 	lis, err := net.Listen("tcp", ":80")
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}

// 	var opts []googleGrpc.ServerOption
// 	grpcServer := googleGrpc.NewServer(opts...)
// 	proto.RegisterToolsServer(grpcServer, &grpc.ToolsServer{})
// 	grpcServer.Serve(lis)

// }
