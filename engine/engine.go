package engine

import (
	"fmt"
	"net"
	"time"

	"github.com/ZyrnDev/letsgohabits/config"
	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/ZyrnDev/letsgohabits/grpc"
	"github.com/ZyrnDev/letsgohabits/nats"
	"github.com/ZyrnDev/letsgohabits/proto"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	googleGrpc "google.golang.org/grpc"
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

type Engine struct {
	natsConneciton *nats.Connection
	database       database.Database
	scheduler      *gocron.Scheduler
	grpc           *googleGrpc.Server
	grpcListener   net.Listener
}

func New(args ...string) (*Engine, error) {
	var engine Engine

	log.Info().Strs("args", args).Msg("Starting Engine")

	conf, err := config.New[ServerConfig]("engine.toml")
	if err != nil {
		return nil, fmt.Errorf("Failed to load config: %s", err)
	} else {
		log.Info().Msgf("Loaded Config: %+v", conf)
	}

	engine.database, err = database.New(conf.DatabaseConfig.ConnectionString, &database.Config{
		// Logger: logger.Default.LogMode(logger.Info), // Verbose Logging
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %s", err)
	}

	engine.natsConneciton, err = nats.Connect(conf.NatsConfig.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to nats: %s", err)
	}
	// defer nc.Close() // TODO: This is causing issue as it closes too early

	log.Info().Msgf("Connected to nats: %s", conf.NatsConfig.ConnectionString)

	engine.grpcListener, err = net.Listen("tcp", ":9090")
	if err != nil {
		return nil, fmt.Errorf("Failed to open tcp listener: %v", err)
	}

	var opts []googleGrpc.ServerOption
	engine.grpc = googleGrpc.NewServer(opts...)
	proto.RegisterToolsServer(engine.grpc, &grpc.ToolsServer{})
	proto.RegisterUsersServer(engine.grpc, &grpc.UsersServer{Database: engine.database})
	proto.RegisterHabitsServer(engine.grpc, &grpc.HabitsServer{Database: engine.database})

	engine.scheduler = gocron.NewScheduler(time.Local)
	go engine.Start()

	return &engine, nil
}

func (engine *Engine) Start() {
	engine.scheduler.Every(1).Second().Do(func() {
		engine.natsConneciton.Publish("print", []byte("Hello World"))
	})
	engine.scheduler.StartAsync()

	engine.grpc.Serve(engine.grpcListener)
}

func (engine *Engine) Close() {
	engine.natsConneciton.Close()
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
// 		log.Fatal().Msgf("failed to listen: %v", err)
// 	}

// 	var opts []googleGrpc.ServerOption
// 	grpcServer := googleGrpc.NewServer(opts...)
// 	proto.RegisterToolsServer(grpcServer, &grpc.ToolsServer{})
// 	grpcServer.Serve(lis)

// }
