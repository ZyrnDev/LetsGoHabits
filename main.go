package main

import (
	"log"
	"os/exec"

	"github.com/ZyrnDev/letsgohabits/client"
	"github.com/ZyrnDev/letsgohabits/config"
	"github.com/ZyrnDev/letsgohabits/database"
	"github.com/ZyrnDev/letsgohabits/server"
)

func main() {
	cmd := exec.Command("cp", "--recursive", "proto/", "generated_proto/")
	cmd.Run()

	conf := config.New()
	log.Printf("Loaded Config: %+v", conf)

	shutdownRequested := make(chan bool)

	db := database.New(conf.DatabaseConnectionString, &database.Config{
		// Logger: logger.Default.LogMode(logger.Info), // Verbose Logging
	})

	shutdownClient := make(chan bool)
	clientDone := client.New(conf.NatsConnectionString, db, shutdownClient)
	shutdownServer := make(chan bool)
	serverDone := server.New(conf.NatsConnectionString, db, shutdownServer)

	// go func() {
	// 	time.Sleep(time.Second * 10)
	// 	log.Println("Shutting down")
	// 	shutdownClient <- true
	// 	shutdownServer <- true
	// 	shutdownRequested <- true
	// }()

	<-clientDone
	<-serverDone
	<-shutdownRequested
}
