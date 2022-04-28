package main

import (
	"time"

	"github.com/ZyrnDev/letsgohabits/handler"
	"github.com/ZyrnDev/letsgohabits/util"
)

func main() {
	shutdownRequests := util.SetupShutdown(util.ShutdownTimeouts{KillTimeout: time.Second * 1, InterruptTimeout: time.Second * 1})

	h, err := handler.New()
	if err != nil {
		panic(err)
	}

	<-shutdownRequests
	h.Close()
}
