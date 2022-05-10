package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/ZyrnDev/letsgohabits/config"
	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/ZyrnDev/letsgohabits/nats"
	"github.com/ZyrnDev/letsgohabits/proto"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	GRPC_TIMEOUT = 10 * time.Second
)

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

var defaultConfig ClientConfig = ClientConfig{
	NatsConfig: NatsConfig{
		ConnectionString: "nats://localhost:4222",
	},
}

type Handler struct {
	natsConnection    *nats.Connection
	grpcConnection    *grpc.ClientConn
	toolsGRPC         proto.ToolsClient
	usersGRPC         proto.UsersClient
	habitsGRPC        proto.HabitsClient
	subscriptionsGRPC proto.SubscriptionsClient
}

func New(args ...string) (*Handler, error) {
	var handler Handler

	conf, err := config.New[ClientConfig](defaultConfig)
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
	handler.toolsGRPC = proto.NewToolsClient(handler.grpcConnection)
	handler.usersGRPC = proto.NewUsersClient(handler.grpcConnection)
	handler.habitsGRPC = proto.NewHabitsClient(handler.grpcConnection)
	handler.subscriptionsGRPC = proto.NewSubscriptionsClient(handler.grpcConnection)

	go handler.SetupGin()

	// handler.natsConnection.Subscribe("print", func(msg nats.NatsMsg) {

	// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 	pingTime, err := handler.grpcClient.Ping(ctx, &empty.Empty{})
	// 	if err != nil {
	// 		log.Info().Err(err).Msgf("Received message '%s' at an unknown time.", msg.Data)
	// 	} else {
	// 		log.Info().Msgf("Received message '%s' at %v", msg.Data, pingTime.AsTime())
	// 	}
	// 	defer cancel()
	// })

	return &handler, nil
}

func (handler *Handler) ToolsPing() (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	pingTime, err := handler.toolsGRPC.Ping(ctx, &empty.Empty{})
	if err != nil {
		return time.Time{}, err
	}
	return pingTime.AsTime(), nil
}

func (handler *Handler) Close() {
	handler.natsConnection.Close()
	handler.grpcConnection.Close()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (handler *Handler) GinPing(c *gin.Context) {
	pingTime, err := handler.ToolsPing()
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "pong",
			"time":    pingTime.Format(time.RFC3339),
		})
	}
}

type GrpcExecutor[Input any] func(input *Input) (interface{}, error)

func GinGrpcCall[Input any](c *gin.Context, grpcOperation GrpcExecutor[Input]) {
	var input Input

	if err := c.BindJSON(&input); err == nil {
		data, err := grpcOperation(&input)

		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "pong",
				"data":    data,
			})
		}
	} else {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("%+v", err),
		})
	}
}

// {"title":"Nice Bro","body":"lets get this bread","vibrate":[100,50,100],"actions":[{"action":"close","title":"Close notification","icon":"https://test.zyrn.dev/gigachad.jpg"}]}
type Notification struct {
	Title   string `json:"title" binding:"required"`
	Body    string `json:"body" binding:"required"`
	Vibrate []int  `json:"vibrate,omitempty"`
	Actions []struct {
		Action string `json:"action" binding:"required"`
		Title  string `json:"title" binding:"required"`
		Icon   string `json:"icon" binding:"required"`
	} `json:"actions,omitempty"`
}

func sendSuccess(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, gin.H{
		"result": "success",
		"data":   data,
	})
}

func sendError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"result": "error",
		"error":  err.Error(),
	})
}

func sendOk(c *gin.Context, data any) {
	sendSuccess(c, 200, data)
}

func sendCreated(c *gin.Context, data any) {
	sendSuccess(c, 201, data)
}

func sendClientError(c *gin.Context, err error) {
	sendError(c, 400, err)
}

func sendNotFound(c *gin.Context, err error) {
	sendError(c, 404, err)
}

func sendServerError(c *gin.Context, err error) {
	sendError(c, 500, err)
}

func (handler *Handler) SetupGin() {

	r := gin.Default()
	r.Use(CORSMiddleware())

	// Ping Service
	r.GET("/ping", handler.GinPing)

	// Users Service
	r.POST("/users/find", handler.GinFindUsers)
	r.POST("/users/create", handler.GinNewUser)
	r.POST("/users/delete", handler.GinDeleteUser)
	r.POST("/users/update", handler.GinUpdateUser)

	// Habits Service
	r.POST("/habits/find", handler.GinFindHabits)
	r.POST("/habits/create", handler.GinNewHabit)
	r.POST("/habits/delete", handler.GinDeleteHabit)
	r.POST("/habits/update", handler.GinUpdateHabit)

	r.POST("/subscriptions/find", handler.GinFindSubscriptions)
	r.POST("/subscriptions/create", handler.GinNewSubscription)
	r.POST("/subscriptions/delete", handler.GinDeleteSubscription)
	r.POST("/subscriptions/update", handler.GinUpdateSubscription)

	r.POST("/users/subscribe", func(c *gin.Context) {
		var input database.Subscription

		if err := c.BindJSON(&input); err != nil {
			sendClientError(c, err)
			return
		}

		subscription := input.ToWebpushSubScription()

		notifcation := Notification{
			Title: "Nice Bro Test",
			Body:  "This is a test notification",
		}
		message, err := json.Marshal(notifcation)
		if err != nil {
			sendServerError(c, err)
			return
		}

		resp, err := webpush.SendNotification(message, &subscription, &webpush.Options{
			Subscriber:      "mitch@zyrn.dev",
			VAPIDPublicKey:  "BH3wNAfqSwzLzduT_KsuE0PoIKRMooJQ_On_iv6uQfNWDc5CNXsH6GYErTQ-OrOTe3LO6H32fJV6eyINpsqHpDg",
			VAPIDPrivateKey: "-DWpYb3--uc2Jo1QRpFAGnZigtBoSgItlwUtP3JN0f8",
			TTL:             30,
		})
		if err != nil {
			sendServerError(c, err)
			return
		}
		defer resp.Body.Close()

		sendOk(c, "nice")
	})
	// r.POST("/users/unsubscribe", func(c *gin.Context) {})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
