package client

import (
	"github.com/ZyrnDev/letsgohabits/config"
	"github.com/ZyrnDev/letsgohabits/nats"
	"github.com/rs/zerolog/log"
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

func New(args ...string) {
	conf, err := config.New[ClientConfig]("client.toml")
	if err != nil {
		log.Fatal().Msgf("Failed to load config: %s", err)
	} else {
		log.Info().Msgf("Loaded Config: %+v", conf)
	}

	log.Info().Msgf("Starting client %+v", args)

	nc, err := nats.Connect(conf.NatsConfig.ConnectionString)
	if err != nil {
		log.Fatal().Msgf("Failed to connect to nats: %s", err)
	}
	// defer nc.Close() // TODO: This is causing issue as it closes too early

	log.Info().Msgf("Connected to nats: %s", conf.NatsConfig.ConnectionString)

	nc.Subscribe("print", func(msg nats.NatsMsg) {
		log.Info().Msgf("Received message: %s", msg.Data)
	})
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

// func handler(w http.ResponseWriter, req *http.Request) {
// 	log.Printf("Received request: %s", req.URL.Path)
// 	fmt.Fprintf(w, "<h1>Hello World</h1><p>This is a test</p>")
// }
