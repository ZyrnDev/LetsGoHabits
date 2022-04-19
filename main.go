package main

import (
	"os/exec"
	"time"

	"github.com/ZyrnDev/letsgohabits/client"
	"github.com/ZyrnDev/letsgohabits/server"
	"github.com/rs/zerolog/log"
)

func CopyGeneratedProtoFilesToMount() {
	cmd := exec.Command("cp", "--recursive", "proto/", "generated_proto/")
	cmd.Run()
}

func main() {
	go CopyGeneratedProtoFilesToMount()

	clientDone := make(chan bool)
	serverDone := make(chan bool)

	go func() { server.New(); time.Sleep(time.Second * 10); serverDone <- true }()
	go func() { client.New(); time.Sleep(time.Second * 10); clientDone <- true }()

	for i := 0; i < 2; i++ {
		select {
		case <-clientDone:
			log.Info().Msg("Client done")
		case <-serverDone:
			log.Info().Msg("Server done")
		}
	}

	log.Info().Msg("Both Client and Server terminated: shutting down")
}
