package main

import (
	"time"

	"github.com/ZyrnDev/letsgohabits/engine"
	"github.com/ZyrnDev/letsgohabits/util"
)

func main() {
	shutdownRequests := util.SetupShutdown(util.ShutdownTimeouts{KillTimeout: time.Second * 1, InterruptTimeout: time.Second * 1})

	e, err := engine.New()
	if err != nil {
		panic(err)
	}

	<-shutdownRequests
	e.Close()
}
