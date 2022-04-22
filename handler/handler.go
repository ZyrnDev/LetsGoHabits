package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ZyrnDev/letsgohabits/config"
	"github.com/ZyrnDev/letsgohabits/mounts/proto/proto"
	"github.com/ZyrnDev/letsgohabits/nats"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// "context"
// "fmt"
// "log"
// "net/http"
// "time"

// "github.com/ZyrnDev/letsgohabits/database"
// "github.com/ZyrnDev/letsgohabits/nats"
// "github.com/ZyrnDev/letsgohabits/proto"
// "github.com/golang/protobuf/ptypes/empty"
// "google.golang.org/grpc"
// "google.golang.org/grpc/credentials/insecure"
// protobuf "google.golang.org/protobuf/proto"
// "gorm.io/gorm"

type NatsConfig struct {
	ConnectionString string `mapstructure:"connection_string"`
}

type ClientConfig struct {
	NatsConfig `mapstructure:"nats"`
}

type Handler struct {
	natsConnection *nats.Connection
	grpcConnection *grpc.ClientConn
	grpcClient     proto.ToolsClient
}

func New(args ...string) (*Handler, error) {
	var handler Handler

	conf, err := config.New[ClientConfig]("handler.toml")
	if err != nil {
		return nil, fmt.Errorf("Failed to load config: %s", err)
	} else {
		log.Info().Msgf("Loaded Config: %+v", conf)
	}

	log.Info().Strs("args", args).Msg("Starting client")

	handler.natsConnection, err = nats.Connect(conf.NatsConfig.ConnectionString)
	if err != nil {
		log.Fatal().Msgf("Failed to connect to nats: %s", err)
	}

	log.Info().Msgf("Connected to nats: %s", conf.NatsConfig.ConnectionString)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	time.Sleep(time.Second * 1)

	handler.grpcConnection, err = grpc.Dial("engine:9090", opts...)
	if err != nil {
		log.Fatal().Msgf("fail to dial: %v", err)
	}
	handler.grpcClient = proto.NewToolsClient(handler.grpcConnection)

	handler.natsConnection.Subscribe("print", func(msg nats.NatsMsg) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pingTime, err := handler.grpcClient.Ping(ctx, &empty.Empty{})
		if err != nil {
			log.Info().Err(err).Msgf("Received message '%s' at an unknown time.", msg.Data)
		} else {
			log.Info().Msgf("Received message '%s' at %v", msg.Data, pingTime.AsTime())
		}
		defer cancel()
	})

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			log.Printf("Received request: %s", req.URL.Path)
			t, err := (&handler).Ping()
			if err != nil {
				fmt.Fprintf(w, "<h1>Hello World</h1><p>Error Ping Failed: %v</p>", err)
			} else {
				fmt.Fprintf(w, "<h1>Hello World</h1><p>Ping Successful: %v</p>", t)
			}
		})
		http.ListenAndServe(":80", nil)
	}()

	return &handler, nil
}

func (handler *Handler) Ping() (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pingTime, err := handler.grpcClient.Ping(ctx, &empty.Empty{})
	if err != nil {
		return time.Time{}, err
	}
	return pingTime.AsTime(), nil
}

func (handler *Handler) Close() {
	handler.natsConnection.Close()
	handler.grpcConnection.Close()
}

// func New(natsConnStr string, db database.Database, shutdownRequested chan bool) chan bool {
// 	done := make(chan bool)

// 	go func() {
// 		nc := nats.NatsConnection(natsConnStr)
// 		defer nc.Close()

// 		log.Printf("Client connected to %s", natsConnStr)

// 		nc.Subscribe("print", func(msg nats.NatsMsg) {
// 			log.Printf("Received: %s", string(msg.Data))
// 		})

// 		nc.Subscribe("test", func(msg nats.NatsMsg) {
// 			t := &proto.Test{}
// 			err := protobuf.Unmarshal(msg.Data, t)
// 			if err != nil {
// 				panic(err)
// 			}

// 			var habit database.Habit
// 			res := db.Preload("Author").First(&habit, database.Habit{Model: gorm.Model{ID: uint(t.Id)}})
// 			if res.Error != nil {
// 				panic(res.Error)
// 			}

// 			log.Printf("%s just send a message (%s) to add a new habit '%s'", habit.Author.Nickname, t.Name, habit.Name)
// 			// log.Printf("Info %+v", habit)
// 		})

// 		// go func() {
// 		// 	http.HandleFunc("/", handler)
// 		// 	http.ListenAndServe(":80", nil)
// 		// }()

// 		go test()

// 		<-shutdownRequested
// 		done <- true
// 	}()

// 	return done
// }

// func test() {
// 	var opts []grpc.DialOption
// 	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

// 	time.Sleep(time.Second * 1)

// 	conn, err := grpc.Dial(":80", opts...)
// 	if err != nil {
// 		log.Fatalf("fail to dial: %v", err)
// 	}
// 	defer conn.Close()

// 	client := proto.NewToolsClient(conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	pingTime, err := client.Ping(ctx, &empty.Empty{})
// 	if err != nil {
// 		log.Printf("fail to ping: %v", err)
// 	} else {
// 		log.Printf("ping time: %v", pingTime.AsTime())
// 	}
// 	defer cancel()
// }
