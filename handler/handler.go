package handler

import (
	"context"
	"fmt"
	"time"

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
	natsConnection *nats.Connection
	grpcConnection *grpc.ClientConn
	toolsGRPC      proto.ToolsClient
	usersGRPC      proto.UsersClient
	habitsGRPC     proto.HabitsClient
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

// TODO(Mitch): Make this more generic using a wrapper function
func (handler *Handler) UsersGet(user *proto.User) (*proto.ListUsers, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	users, err := handler.usersGRPC.Get(ctx, user)
	return users, err
}

func (handler *Handler) UsersNew(user *proto.User) (*proto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	user, err := handler.usersGRPC.New(ctx, user)
	return user, err
}

func (handler *Handler) UsersDelete(user *proto.User) (*empty.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	empty_, err := handler.usersGRPC.Delete(ctx, user)
	return empty_, err
}

func (handler *Handler) UserUpdate(user *proto.User) (*proto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	user, err := handler.usersGRPC.Update(ctx, user)
	return user, err
}

func (handler *Handler) HabitsGet(habit *proto.Habit) (*proto.ListHabits, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	habits, err := handler.habitsGRPC.Get(ctx, habit)
	return habits, err
}

func (handler *Handler) HabitsNew(habit *proto.Habit) (*proto.Habit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	habit, err := handler.habitsGRPC.New(ctx, habit)
	return habit, err
}

func (handler *Handler) HabitsDelete(habit *proto.Habit) (*empty.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	empty_, err := handler.habitsGRPC.Delete(ctx, habit)
	return empty_, err
}

func (handler *Handler) HabitUpdate(habit *proto.Habit) (*proto.Habit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPC_TIMEOUT)
	defer cancel()
	habit, err := handler.habitsGRPC.Update(ctx, habit)
	return habit, err
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

func (handler *Handler) GinFindUsers(c *gin.Context) {
	GinGrpcCall(c, func(user *database.User) (interface{}, error) {
		return handler.UsersGet(user.ToProtobuf())
	})
}

func (handler *Handler) GinNewUser(c *gin.Context) {
	GinGrpcCall(c, func(user *database.User) (interface{}, error) {
		if user.ID != 0 {
			return nil, fmt.Errorf("user.ID must be not set or 0")
		}
		if user.Nickname == "" {
			return nil, fmt.Errorf("user.Nickname must be set and not empty")
		}
		return handler.UsersNew(user.ToProtobuf())
	})
}

func (handler *Handler) GinDeleteUser(c *gin.Context) {
	GinGrpcCall(c, func(user *database.User) (interface{}, error) {
		if user.ID == 0 {
			return nil, fmt.Errorf("user.ID must be set and not empty")
		}
		return handler.UsersDelete(user.ToProtobuf())
	})
}

func (handler *Handler) GinUpdateUser(c *gin.Context) {
	GinGrpcCall(c, func(user *database.User) (interface{}, error) {
		if user.ID == 0 {
			return nil, fmt.Errorf("user.ID must be set and not empty")
		}
		return handler.UserUpdate(user.ToProtobuf())
	})
}

func (handler *Handler) GinFindHabits(c *gin.Context) {
	GinGrpcCall(c, func(habit *database.Habit) (interface{}, error) {
		return handler.HabitsGet(habit.ToProtobuf())
	})
}

func (handler *Handler) GinNewHabit(c *gin.Context) {
	GinGrpcCall(c, func(habit *database.Habit) (interface{}, error) {
		if habit.ID != 0 {
			return nil, fmt.Errorf("habit.ID must be not set or 0")
		}
		if habit.Name == "" {
			return nil, fmt.Errorf("habit.Name must be set and not empty")
		}
		return handler.HabitsNew(habit.ToProtobuf())
	})
}

func (handler *Handler) GinDeleteHabit(c *gin.Context) {
	GinGrpcCall(c, func(habit *database.Habit) (interface{}, error) {
		if habit.ID == 0 {
			return nil, fmt.Errorf("habit.ID must be set and not empty")
		}
		return handler.HabitsDelete(habit.ToProtobuf())
	})
}

func (handler *Handler) GinUpdateHabit(c *gin.Context) {
	GinGrpcCall(c, func(habit *database.Habit) (interface{}, error) {
		if habit.ID == 0 {
			return nil, fmt.Errorf("habit.ID must be set and not empty")
		}
		return handler.HabitUpdate(habit.ToProtobuf())
	})
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

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
